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

//今日导读

const CALl_BLOCK_TODAY_INTRO = "call_block_today_introduction"
const CALl_BLOCK_TODAY_INTRO_TRAIT = "call_block_today_introduction_trait"

type IntroRules struct {
	Day     int `json:"day"`
	CronExp int `json:"cronExp"` // 周期时间
	Limit      int `json:"limit"`
}

type TodayIntro struct {
	topicId int
	name    string

	introRules IntroRules

	cancel context.CancelFunc
	Ctx    context.Context

	topicIds []string // 数据源
}

func NewTodayIntro(topicId int, topicIds []string) *TodayIntro {
	return &TodayIntro{
		topicId:  topicId,
		name:     fmt.Sprintf("%d%s%s", topicId, service.Separator, service.TODAY_INTRO),
		topicIds: topicIds,
	}
}

// 删除 redis 数据
func (this *TodayIntro) Remover(tids []int) {
	this.remover(tids)
}

func (this *TodayIntro) remover(tids []int) {
	data_source.DelRedisThreadInfo(tids, this.redisKey(), this.traitRedisKey())
}

func (this *TodayIntro) GetThis() interface{} {
	return this
}

func (this *TodayIntro) Init() {
	ctx, cancel := context.WithCancel(context.Background())
	this.Ctx = ctx
	this.cancel = cancel
	this.introRules = todayIntro
}

func (this *TodayIntro) ChangeConf(conf string) error {
	var rule IntroRules
	if err := json.Unmarshal([]byte(conf), &rule); err == nil {
		this.introRules = rule
		go this.reStart()
		return nil
	} else {
		return err
	}
}

func (this *TodayIntro) Start() {
	this.worker()
	t := time.NewTimer(time.Duration(this.introRules.CronExp) * time.Minute)
	for {
		select {
		case <-t.C:
			this.worker()
			t.Reset(time.Duration(this.introRules.CronExp) * time.Minute)
		case <-this.Ctx.Done():
			return
		}
	}
}

// 写 reids
func (this *TodayIntro) worker() {
	redisThreads := data_source.GetTodayIntroSortThread(topic_fid_relation.GetFids(this.topicIds), this.introRules.Day,this.introRules.Limit)
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

func (this *TodayIntro) ChangeFids(topicIds []string) {
	this.topicIds = topicIds
	go this.reStart()
}

func (this *TodayIntro) Stop() {
	boot.InstanceRedisCli(boot.CACHE).Del(this.redisKey())
	this.cancel()
}

func (this *TodayIntro) GetName() string {
	return this.name
}

func (this *TodayIntro) reStart() {
	this.Stop()
	this.Ctx, this.cancel = context.WithCancel(context.Background())
	this.Start()
}

func (this *TodayIntro) AddTrait(id string, trait service.CallBlockTrait) {
	if traitBytes, err := json.Marshal(&trait); err == nil {
		boot.InstanceRedisCli(boot.CACHE).HSet(this.traitRedisKey(), id, string(traitBytes))
	}
}

func (this *TodayIntro) redisKey() string {
	return fmt.Sprintf("%s%s%d", CALl_BLOCK_TODAY_INTRO, service.Separator, this.topicId)
}

func (this *TodayIntro) traitRedisKey() string {
	return fmt.Sprintf("%s%s%d", CALl_BLOCK_TODAY_INTRO_TRAIT, service.Separator, this.topicId)
}
