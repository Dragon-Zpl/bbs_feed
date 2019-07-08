package call_block

/* 精华调用块*/

import (
	"context"
	"fmt"
	"time"
)

type EssenceRules struct {
	cronExp time.Duration // 周期时间
	criticalTime time.Duration // 临界时间
}


type Essence struct {
	name string

	reportChan chan []string
	essenceRules EssenceRules

	cancel       context.CancelFunc
	Ctx          context.Context

	topicIds []string // 数据源

}

func NewEssence(topicId int, topicIds []string) *Essence {
	return &Essence{
		name:fmt.Sprintf("%d-essence", topicId),
		topicIds:topicIds,
	}
}

func(this *Essence) RemoveReportThread() {
	for {
		select {
		case tids := <- this.reportChan:
			// todo clear redis uids
			fmt.Println(tids)
		}
	}
}


func (this *Essence) AcceptSign(tids []string){
	this.reportChan <- tids
	return
}

func (this *Essence) GetThis() interface{} {
	return this
}


func(this *Essence) Init() {
	ctx, cancel := context.WithCancel(context.Background())
	this.Ctx = ctx
	this.cancel = cancel
	this.essenceRules = essence
}

func (this *Essence) ChangeConf(conf interface{}) {
	if conf, ok := conf.(EssenceRules); ok {
		this.essenceRules = conf
		this.reStart()
	}
}

func (this *Essence) Start() {
	t := time.NewTimer(this.essenceRules.cronExp)
	for {
		select {
		case <- t.C:
			// todo  根据配置 写redis数据
			t.Reset(this.essenceRules.cronExp)
		case <- this.Ctx.Done():
			return
		}
	}
}


func (this *Essence) Stop() {
	this.cancel()
}

func (this *Essence) ChangeFids(topicIds []string) {
	this.topicIds = topicIds
}

func (this *Essence) GetName() (string){
	return this.name
}

func (this *Essence) reStart() {
	this.Stop()
	this.Ctx, this.cancel = context.WithCancel(context.Background())
	this.Start()
}