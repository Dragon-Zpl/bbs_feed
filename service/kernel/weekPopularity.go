package kernel

import (
	"bbs_feed/service/service_confs"
	"context"
	"fmt"
	"time"
)

type WeekPopularityRule struct {
	cronExp time.Duration // 周期时间
}

type WeekPopularity struct {
	name string

	reportChan chan []string
	weekPopularityRule WeekPopularityRule

	cancel       context.CancelFunc
	Ctx          context.Context

	topicIds []string // 数据源

}

func NewWeekPopularity(topicId int, topicIds []string) *WeekPopularity {
	return &WeekPopularity{
		name:fmt.Sprintf("%d-weekPopularity", topicId),
		topicIds:topicIds,
	}
}


func(this *WeekPopularity) RemoveReportUser() {
	for {
		select {
		case uids := <- this.reportChan:
			// todo clear redis uids
		}
	}
}

func (this *WeekPopularity) AcceptSign(userIds []string){
	this.reportChan <- userIds
	return
}

func (this *WeekPopularity) GetThis() interface{} {
	return this
}


func(this *WeekPopularity) Init() {
	ctx, cancel := context.WithCancel(context.Background())
	this.Ctx = ctx
	this.cancel = cancel
	this.weekPopularityRule = service_confs.WeekPopularity
}

func (this *WeekPopularity) ChangeConf(conf interface{}) {
	if conf, ok := conf.(WeekPopularityRule); ok {
		this.weekPopularityRule = conf
		this.reStart()
	}
}

func (this *WeekPopularity) Start() {
	t := time.NewTimer(this.weekPopularityRule.cronExp)
	for {
		select {
		case <- t.C:
			// todo  根据配置 写redis数据
			t.Reset(this.weekPopularityRule.cronExp)
		case <- this.Ctx.Done():
			return
		}
	}
}

func (this *WeekPopularity) ChangeFids(topicIds []string) {
	this.topicIds = topicIds
}

func (this *WeekPopularity) Stop() {
	this.cancel()
}


func (this *WeekPopularity) GetName() (string){
	return this.name
}

func (this *WeekPopularity) reStart() {
	this.Stop()
	this.Ctx, this.cancel = context.WithCancel(context.Background())
	this.Start()
}