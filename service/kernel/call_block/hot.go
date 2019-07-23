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

//热帖

const CALl_BLOCK_HOT_THREAD = "call_block_hot"
const CALl_BLOCK_HOT_THREAD_TRAIT = "call_block_hot_trait"

type HotRules struct {
	Day        int `json:"day"` //几天内的数据
	ViewCount  int `json:"viewCount"`
	ReplyCount int `json:"replyCount"`
	CronExp    int `json:"cronExp"` // 周期时间
	Limit      int `json:"limit"`
}

type Hot struct {
	topicId int
	name    string

	hotRules HotRules

	cancel context.CancelFunc
	Ctx    context.Context

	topicIds []string // 数据源

	contract.ThreadRep
}

func NewHot(topicId int, topicIds []string) *Hot {
	hot := &Hot{
		topicId:  topicId,
		name:     fmt.Sprintf("%d%s%s", topicId, service.Separator, service.HOT),
		topicIds: topicIds,
	}
	hot.ReportChan = make(chan []int, 10)
	return hot
}

// 删除 redis 数据
func (this *Hot) Remover(tids []int) {
	this.remover(tids)
}

func (this *Hot) remover(tids []int) {
	logs.Info("remove --", this.redisKey(), "--", this.traitRedisKey(), "--", tids)
	data_source.DelRedisThreadInfo(tids, this.redisKey(), this.traitRedisKey())
}

func (this *Hot) GetThis() interface{} {
	return this
}

func (this *Hot) Init() {
	ctx, cancel := context.WithCancel(context.Background())
	this.Ctx = ctx
	this.cancel = cancel
	this.hotRules = hot
	go this.RemoveReportThread(this.remover) // 开启举报帖自检
}

func (this *Hot) ChangeConf(conf string) error {
	var rule HotRules
	if err := json.Unmarshal([]byte(conf), &rule); err == nil {
		this.hotRules = rule
		go this.reStart()
		return nil
	} else {
		return err
	}
}

func (this *Hot) Start() {
	this.worker()
	t := time.NewTimer(time.Duration(this.hotRules.CronExp) * time.Minute)
	for {
		select {
		case <-t.C:
			this.worker()
			t.Reset(time.Duration(this.hotRules.CronExp) * time.Minute)
		case <-this.Ctx.Done():
			return
		}
	}
}

// 写 reids
func (this *Hot) worker() {
	redisThreads := data_source.GetHotSortThread(topic_fid_relation.GetFids(this.topicIds), this.hotRules.Day, this.hotRules.ViewCount, this.hotRules.ReplyCount, this.hotRules.Limit)
	redisTraits, _ := boot.InstanceRedisCli(boot.CACHE).HGetAll(this.traitRedisKey()).Result()

	datas := make([]interface{}, 0, len(redisThreads))
	for _, thread := range redisThreads {
		if redisTraits != nil {
			if threadTrait, ok := redisTraits[strconv.Itoa(thread.Thread.Tid)]; ok {
				var callBlockTrait service.CallBlockTrait
				if err := json.Unmarshal([]byte(threadTrait), &callBlockTrait); err == nil {
					expTime, _ := time.Parse("2006-01-02 15:04:05", callBlockTrait.Exp)
					if time.Now().Sub(expTime) < 0 {
						thread.Trait = callBlockTrait
					} else {
						redis_ops.Hdel(this.traitRedisKey(), strconv.Itoa(thread.Thread.Tid))
					}
				}
			}
		}

		datas = append(datas, thread)
	}
	redis_ops.ZAddSort(this.redisKey(), datas)
	logs.Info(this.redisKey(), "insert success")
}

func (this *Hot) ChangeFids(topicIds []string) {
	this.topicIds = topicIds
	go this.reStart()
}

func (this *Hot) Stop() {
	boot.InstanceRedisCli(boot.CACHE).Del(this.redisKey())
	logs.Info(this.redisKey(), "delete success")
	this.cancel()
}

func (this *Hot) GetName() string {
	return this.name
}

func (this *Hot) reStart() {
	this.Stop()
	this.Ctx, this.cancel = context.WithCancel(context.Background())
	this.Start()
}

func (this *Hot) AddTrait(id string, trait service.CallBlockTrait) {
	if traitBytes, err := json.Marshal(&trait); err == nil {
		boot.InstanceRedisCli(boot.CACHE).HSet(this.traitRedisKey(), id, string(traitBytes))
	}
}

func (this *Hot) redisKey() string {
	return fmt.Sprintf("%s%s%d", CALl_BLOCK_HOT_THREAD, service.Separator, this.topicId)
}

func (this *Hot) traitRedisKey() string {
	return fmt.Sprintf("%s%s%d", CALl_BLOCK_HOT_THREAD_TRAIT, service.Separator, this.topicId)
}
