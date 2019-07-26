package call_block

import (
	"bbs_feed/boot"
	"bbs_feed/service"
	"bbs_feed/service/data_source"
	"bbs_feed/service/kernel/contract"
	"bbs_feed/service/redis_ops"
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"strconv"
	"time"
)

//周人气榜

const CALl_BLOCK_WEEK_POPULARITY = "call_block_week_popularity"
const CALl_BLOCK_WEEK_POPULARITY_TRAIT = "call_block_week_popularity_trait"

type WeekPopularityRule struct {
	CronExp time.Duration `json:"cronExp"` // 周期时间
	Limit   int           `json:"limit"`
}

type WeekPopularity struct {
	topicId int
	name    string

	weekPopularityRule WeekPopularityRule

	cancel context.CancelFunc
	Ctx    context.Context

	topicIds []string // 数据源

	contract.UserRep //违规用户
}

func NewWeekPopularity() *WeekPopularity {
	return &WeekPopularity{}
}

// 删除 redis 数据
func (this *WeekPopularity) Remover(uids []int) {
	this.remover(uids)
}

func (this *WeekPopularity) remover(uids []int) {
	logs.Info("remove --", this.redisKey(), "--", this.traitRedisKey(), "--", uids)
	data_source.DelRedisThreadInfo(uids, this.redisKey(), this.traitRedisKey())
}

func (this *WeekPopularity) GetThis() interface{} {
	return this
}

func (this *WeekPopularity) Init() {
	ctx, cancel := context.WithCancel(context.Background())
	this.Ctx = ctx
	this.cancel = cancel
	this.weekPopularityRule = weekPopularity
	go this.RemoveReportUser(this.remover) // 开启违规用户自检
}

func (this *WeekPopularity) ChangeConf(conf string) error {
	var rule WeekPopularityRule
	if err := json.Unmarshal([]byte(conf), &rule); err == nil {
		this.weekPopularityRule = rule
		go this.reStart()
		return nil
	} else {
		return err
	}
}

func (this *WeekPopularity) Start() {
	this.worker()
	t := time.NewTimer(time.Duration(this.weekPopularityRule.CronExp) * time.Minute)
	for {
		select {
		case <-t.C:
			this.worker()
			t.Reset(time.Duration(this.weekPopularityRule.CronExp) * time.Minute)
		case <-this.Ctx.Done():
			return
		}
	}
}

func (this *WeekPopularity) worker() {
	data_result := data_source.GetPopulData(this.weekPopularityRule.Limit)

	redisTraits, _ := boot.InstanceRedisCli(boot.CACHE).HGetAll(this.traitRedisKey()).Result()

	for k, v := range data_result {
		topic_id, _ := strconv.Atoi(k)
		this.topicId = topic_id
		datas := make([]interface{}, 0, len(v))
		for _, user := range v {
			if redisTraits != nil {
				if threadTrait, ok := redisTraits[strconv.Itoa(user.User.Uid)]; ok {
					var callBlockTrait service.CallBlockTrait
					if err := json.Unmarshal([]byte(threadTrait), &callBlockTrait); err == nil {
						expTime, _ := time.Parse("2006-01-02 15:04:05", callBlockTrait.Exp)
						if time.Now().Sub(expTime) < 0 {
							user.Trait = callBlockTrait
						} else {
							redis_ops.Hdel(this.traitRedisKey(), strconv.Itoa(user.User.Uid))
						}
					}
				}
			}
			datas = append(datas, user)
		}
		redis_ops.ZAddSort(this.redisKey(), datas)
	}
}

func (this *WeekPopularity) ChangeFids(topicIds []string) {
	this.topicIds = topicIds
	go this.reStart()
}

func (this *WeekPopularity) Stop() {
	boot.InstanceRedisCli(boot.CACHE).Del(this.redisKey())
	logs.Info(this.redisKey(), "delete success")
	this.cancel()
}

func (this *WeekPopularity) GetName() string {
	return this.name
}

func (this *WeekPopularity) reStart() {
	this.Stop()
	this.Ctx, this.cancel = context.WithCancel(context.Background())
	this.Start()
}

func (this *WeekPopularity) AddTrait(id string, trait service.CallBlockTrait) {
	if traitBytes, err := json.Marshal(&trait); err == nil {
		boot.InstanceRedisCli(boot.CACHE).HSet(this.traitRedisKey(), id, string(traitBytes))
	}
}
func (this *WeekPopularity) redisKey() string {
	return fmt.Sprintf("%s%s%d", CALl_BLOCK_WEEK_POPULARITY, service.Separator, this.topicId)
}

func (this *WeekPopularity) traitRedisKey() string {
	return fmt.Sprintf("%s%s%d", CALl_BLOCK_WEEK_POPULARITY_TRAIT, service.Separator, this.topicId)
}
