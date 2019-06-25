package kernel

import (
	"bbs_feed/service/service_confs"
	"context"
	"fmt"
	"strings"
	"time"
)

type HotRules struct {
	viewCount int
	replyCount int
	cronExp time.Duration // 周期时间
}

type Hot struct {
	name string

	reportChan chan []string
	hotRules HotRules

	cancel       context.CancelFunc
	Ctx          context.Context

	topicIds []string // 数据源

}

func NewHot(topicId int, topicIds []string) *Hot {
	return &Hot{
		name:fmt.Sprintf("%d-hot", topicId),
		topicIds:topicIds,
	}
}



func(this *Hot) RemoveReportThread() {
	for {
		select {
		case tids := <- this.reportChan:
			// todo clear redis uids
		}
	}
}


func (this *Hot) AcceptSign(tids []string){
	this.reportChan <- tids
	return
}

func (this *Hot) GetThis() interface{} {
	return this
}

func(this *Hot) Init() {
	ctx, cancel := context.WithCancel(context.Background())
	this.Ctx = ctx
	this.cancel = cancel
	this.hotRules = service_confs.Hot
}

func (this *Hot) ChangeConf(conf interface{}) {
	if conf, ok := conf.(HotRules); ok {
		this.hotRules = conf
		this.reStart()
	}
}

func (this *Hot) Start() {
	t := time.NewTimer(this.hotRules.cronExp)
	for {
		select {
		case <- t.C:
			// todo  根据配置 写redis数据
			t.Reset(this.hotRules.cronExp)
		case <- this.Ctx.Done():
			return
		}
	}
}

func (this *Hot) ChangeFids(topicIds []string) {
	this.topicIds = topicIds
}


func (this *Hot) Stop() {
	this.cancel()
}



func (this *Hot) GetName() (string){
	return this.name
}

func (this *Hot) reStart() {
	this.Stop()
	this.Ctx, this.cancel = context.WithCancel(context.Background())
	this.Start()
}