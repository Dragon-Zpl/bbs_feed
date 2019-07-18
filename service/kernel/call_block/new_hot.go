package call_block

import (
	"bbs_feed/boot"
	"bbs_feed/model/topic_fid_relation"
	"bbs_feed/service"
	"bbs_feed/service/data_source"
	"bbs_feed/service/redis_ops"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

//最新最热(新闻)

const CALl_BLOCK_NEW_HOT = "call_block_new_hot"
const CALl_BLOCK_NEW_HOT_TRAIT = "call_block_new_hot_trait"

type NewHotRules struct {
	Day     int `json:"day"`
	CronExp int `json:"cronExp"` // 周期时间
	Limit      int `json:"limit"`
}

type NewHots struct {
	topicId int
	name    string

	newHotRules NewHotRules

	cancel context.CancelFunc
	Ctx    context.Context

	topicIds []string // 数据源
}

func NewNewHots(topicId int, topicIds []string) *NewHots {
	return &NewHots{
		topicId:  topicId,
		name:     fmt.Sprintf("%d%s%s", topicId, service.Separator, service.TODAY_INTRO),
		topicIds: topicIds,
	}
}

// 删除 redis 数据
func (this *NewHots) Remover(tids []int) {
	this.remover(tids)
}

func (this *NewHots) remover(tids []int) {
	data_source.DelRedisThreadInfo(tids, this.redisKey(), this.traitRedisKey())
}

func (this *NewHots) GetThis() interface{} {
	return this
}

func (this *NewHots) Init() {
	ctx, cancel := context.WithCancel(context.Background())
	this.Ctx = ctx
	this.cancel = cancel
	this.newHotRules = newHot
}

func (this *NewHots) ChangeConf(conf string) error {
	var rule NewHotRules
	if err := json.Unmarshal([]byte(conf), &rule); err == nil {
		this.newHotRules = rule
		go this.reStart()
		return nil
	} else {
		return err
	}
}

func (this *NewHots) Start() {
	this.worker()
	t := time.NewTimer(time.Duration(this.newHotRules.CronExp) * time.Minute)
	for {
		select {
		case <-t.C:
			this.worker()
			t.Reset(time.Duration(this.newHotRules.CronExp) * time.Minute)
		case <-this.Ctx.Done():
			return
		}
	}
}

// 写 reids
func (this *NewHots) worker() {
	redisThreads := data_source.GetNewHotSortThread(topic_fid_relation.GetFids(this.topicIds), this.newHotRules.Day,this.newHotRules.Limit)
	redisTraits, _ := boot.InstanceRedisCli(boot.CACHE).HGetAll(this.traitRedisKey()).Result()

	datas := make([]interface{}, 0, len(redisThreads))
	for _, thread := range redisThreads {
		if redisTraits != nil {
			if threadTrait, ok := redisTraits[strconv.Itoa(thread.Thread.Tid)]; ok {
				var callBlockTrait service.CallBlockTrait
				if err := json.Unmarshal([]byte(threadTrait), &callBlockTrait); err == nil {
					thread.Trait = callBlockTrait
				}
			}
		}

		datas = append(datas, thread)
	}
	redis_ops.ZAddSort(this.redisKey(), datas)
}

func (this *NewHots) ChangeFids(topicIds []string) {
	this.topicIds = topicIds
	go this.reStart()
}

func (this *NewHots) Stop() {
	boot.InstanceRedisCli(boot.CACHE).Del(this.redisKey())
	this.cancel()
}

func (this *NewHots) GetName() string {
	return this.name
}

func (this *NewHots) reStart() {
	this.Stop()
	this.Ctx, this.cancel = context.WithCancel(context.Background())
	this.Start()
}

func (this *NewHots) AddTrait(id string, trait service.CallBlockTrait) {
	if traitBytes, err := json.Marshal(&trait); err == nil {
		boot.InstanceRedisCli(boot.CACHE).HSet(this.traitRedisKey(), id, string(traitBytes))
	}
}

func (this *NewHots) redisKey() string {
	return fmt.Sprintf("%s%s%d", CALl_BLOCK_NEW_HOT, service.Separator, this.topicId)
}

func (this *NewHots) traitRedisKey() string {
	return fmt.Sprintf("%s%s%d", CALl_BLOCK_NEW_HOT_TRAIT, service.Separator, this.topicId)
}
