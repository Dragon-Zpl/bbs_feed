package call_block

import (
	"bbs_feed/boot"
	"bbs_feed/model/topic_fid_relation"
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
	ThreadSupportedScore int `json:"threadSupportedScore"` //帖子被加分权重
	PublishThreadScore   int `json:"publishThreadScore"`   //发帖权重
	PublishPostScore     int `json:"publishPostScore"`     //评论权重
	PostSupportedScore   int `json:"postSupportScore"`     //评论被加分权重
	CronExp              int `json:"cronExp"`              // 周期时间
	Limit                int `json:"limit"`
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

func NewWeekPopularity(topicId int, topicIds []string) *WeekPopularity {
	weekPopularity := &WeekPopularity{
		topicId:  topicId,
		name:     fmt.Sprintf("%d%s%s", topicId, service.Separator, service.WEEK_POPULARITY),
		topicIds: topicIds,
	}
	weekPopularity.ReportChan = make(chan []int, 10)
	return weekPopularity
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
	data_result := data_source.GetPopulData(topic_fid_relation.GetFids(this.topicIds), this.weekPopularityRule.Limit, this.weekPopularityRule.PostSupportedScore, this.weekPopularityRule.PublishPostScore, this.weekPopularityRule.PublishThreadScore, this.weekPopularityRule.ThreadSupportedScore)

	redisTraits, _ := boot.InstanceRedisCli(boot.CACHE).HGetAll(this.traitRedisKey()).Result()
	datas := make([]interface{}, 0, len(data_result))
	for _, user := range data_result {
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
	logs.Info(this.redisKey(), "insert success")
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
