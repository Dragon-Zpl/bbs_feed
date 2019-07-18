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

//精华贴

const CALl_BLOCK_ESSENCE_THREAD = "call_block_essence"
const CALl_BLOCK_ESSENCE_THREAD_TRAIT = "call_block_essence_trait"

type EssenceRules struct {
	CronExp int `json:"cronExp"` // 周期时间 min
	Day     int `json:"day"`
	Limit   int `json:"limit"`
}

type Essence struct {
	name         string
	topicId      int
	essenceRules EssenceRules

	cancel context.CancelFunc
	Ctx    context.Context

	topicIds []string // 数据源

	contract.ThreadRep
}

func NewEssence(topicId int, topicIds []string) *Essence {
	essence := &Essence{
		topicId:  topicId,
		name:     fmt.Sprintf("%d%s%s", topicId, service.Separator, service.ESSENCE),
		topicIds: topicIds,
	}
	essence.ReportChan = make(chan []int, 10)
	return essence
}

func (this *Essence) remover(tids []int) {
	logs.Info("remove --", this.redisKey(), "--", this.traitRedisKey(), "--", tids)
	data_source.DelRedisThreadInfo(tids, this.redisKey(), this.traitRedisKey())
}

func (this *Essence) GetThis() interface{} {
	return this
}

func (this *Essence) Init() {
	ctx, cancel := context.WithCancel(context.Background())
	this.Ctx = ctx
	this.cancel = cancel
	this.essenceRules = essence
	go this.RemoveReportThread(this.remover)
}

func (this *Essence) ChangeConf(conf string) error {
	var rule EssenceRules
	if err := json.Unmarshal([]byte(conf), &rule); err == nil {
		this.essenceRules = rule
		go this.reStart()
		return nil
	} else {
		return err
	}
}

func (this *Essence) Start() {
	this.worker()
	t := time.NewTimer(time.Duration(this.essenceRules.CronExp) * time.Minute)
	for {
		select {
		case <-t.C:
			this.worker()
			t.Reset(time.Duration(this.essenceRules.CronExp) * time.Minute)
		case <-this.Ctx.Done():
			return
		}
	}
}

// 删除 redis 数据
func (this *Essence) Remover(tids []int) {
	this.remover(tids)
}

func (this *Essence) Stop() {
	boot.InstanceRedisCli(boot.CACHE).Del(this.redisKey())
	this.cancel()
}

func (this *Essence) ChangeFids(topicIds []string) {
	this.topicIds = topicIds
}

func (this *Essence) AddTrait(id string, trait service.CallBlockTrait) {
	if traitBytes, err := json.Marshal(&trait); err == nil {
		boot.InstanceRedisCli(boot.CACHE).HSet(this.traitRedisKey(), id, string(traitBytes))
	}
}

// 写 reids
func (this *Essence) worker() {
	redisThreads := data_source.GetEssenceSortThread(topic_fid_relation.GetFids(this.topicIds), this.essenceRules.Day, this.essenceRules.Limit)
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
	logs.Info(this.redisKey(), "insert success")
}

func (this *Essence) GetName() string {
	return this.name
}

func (this *Essence) reStart() {
	this.Stop()
	this.Ctx, this.cancel = context.WithCancel(context.Background())
	this.Start()
}

func (this *Essence) redisKey() string {
	return fmt.Sprintf("%s%s%d", CALl_BLOCK_ESSENCE_THREAD, service.Separator, this.topicId)
}

func (this *Essence) traitRedisKey() string {
	return fmt.Sprintf("%s%s%d", CALl_BLOCK_ESSENCE_THREAD_TRAIT, service.Separator, this.topicId)
}
