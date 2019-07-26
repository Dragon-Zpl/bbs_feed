package call_block

import (
	"bbs_feed/boot"
	"bbs_feed/service"
	"bbs_feed/service/data_source"
	"bbs_feed/service/kernel/contract"
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"
)

//贡献榜
const CALl_BLOCK_WEEK_CONTRIBUTION = "call_block_week_contribution"
const CALl_BLOCK_WEEK_CONTRIBUTION_TRAIT = "call_block_week_contribution_trait"

type ContributionRules struct {
	Day     int `json:"day"`
	CronExp int `json:"cronExp"` // 周期时间
	Limit   int `json:"limit"`
}

type Contribution struct {
	topicId int
	name    string

	contributionRules ContributionRules

	cancel context.CancelFunc
	Ctx    context.Context

	topicIds []string // 数据源

	contract.UserRep //违规用户
}

func NewContribution(topicId int, topicIds []string) *Contribution {
	return &Contribution{
		name:     fmt.Sprintf("%d%s%s", topicId, service.Separator, service.CONTRIBUTION),
		topicIds: topicIds,
	}
}

// 删除 redis 数据
func (this *Contribution) Remover(uids []int) {
	this.remover(uids)
}

func (this *Contribution) remover(uids []int) {
	logs.Info("remove --", this.redisKey(), "--", this.traitRedisKey(), "--", uids)
	data_source.DelRedisThreadInfo(uids, this.redisKey(), this.traitRedisKey())
}

func (this *Contribution) GetThis() interface{} {
	return this
}

func (this *Contribution) Init() {
	ctx, cancel := context.WithCancel(context.Background())
	this.Ctx = ctx
	this.cancel = cancel
	this.contributionRules = contribution
	go this.RemoveReportUser(this.remover) // 开启违规用户自检
}

func (this *Contribution) ChangeConf(conf string) error {
	var rule ContributionRules
	if err := json.Unmarshal([]byte(conf), &rule); err == nil {
		this.contributionRules = rule
		go this.reStart()
		return nil
	} else {
		return err
	}
}

func (this *Contribution) Start() {
	this.worker()
	t := time.NewTimer(time.Duration(this.contributionRules.CronExp) * time.Minute)
	for {
		select {
		case <-t.C:
			this.worker()
			t.Reset(time.Duration(this.contributionRules.CronExp) * time.Minute)
		case <-this.Ctx.Done():
			return
		}
	}
}

func (this *Contribution) worker() {
	//todo
}

func (this *Contribution) ChangeFids(topicIds []string) {
	this.topicIds = topicIds
	go this.reStart()
}

func (this *Contribution) Stop() {
	boot.InstanceRedisCli(boot.CACHE).Del(this.redisKey())
	logs.Info(this.redisKey(), "delete success")
	this.cancel()
}

func (this *Contribution) GetName() string {
	return this.name
}

func (this *Contribution) reStart() {
	this.Stop()
	this.Ctx, this.cancel = context.WithCancel(context.Background())
	this.Start()
}

func (this *Contribution) AddTrait(id string, trait service.CallBlockTrait) {
	if traitBytes, err := json.Marshal(&trait); err == nil {
		boot.InstanceRedisCli(boot.CACHE).HSet(this.traitRedisKey(), id, string(traitBytes))
	}
}
func (this *Contribution) redisKey() string {
	return fmt.Sprintf("%s%s%d", CALl_BLOCK_WEEK_CONTRIBUTION, service.Separator, this.topicId)
}

func (this *Contribution) traitRedisKey() string {
	return fmt.Sprintf("%s%s%d", CALl_BLOCK_WEEK_CONTRIBUTION_TRAIT, service.Separator, this.topicId)
}
