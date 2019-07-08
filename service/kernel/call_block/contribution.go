package call_block

/*
	贡献榜
*/
import (
	"context"
	"fmt"
	"time"
)


type ContributionRules struct {
	cronExp time.Duration // 周期时间
	views int
	replys int
}

type Contribution struct {
	name string
	reportChan chan []string
	ContributionRules ContributionRules

	cancel       context.CancelFunc
	Ctx          context.Context
	topicIds []string // 数据源
}

func NewContribution(topicId int, topicIds []string) *Contribution {
	return &Contribution{
		name:fmt.Sprintf("%d-contribution", topicId),
		topicIds:topicIds,
	}
}

func(this *Contribution) RemoveReportUser() {
	for {
		select {
		case uids := <- this.reportChan:
			// todo clear redis uids
			fmt.Println(uids)
		}
	}
}


func (this *Contribution) AcceptSign(userIds []string){
	this.reportChan <- userIds
	return
}

func (this *Contribution) GetThis() interface{} {
	return this
}

func (this *Contribution) Init() {
	ctx, cancel := context.WithCancel(context.Background())
	this.Ctx = ctx
	this.cancel = cancel
	this.ContributionRules = contribution
}

func (this *Contribution) ChangeConf(conf interface{}) {
	if conf, ok := conf.(ContributionRules); ok {
		this.ContributionRules = conf
		this.reStart()
	}
}

func (this *Contribution) Start() {
	t := time.NewTimer(this.ContributionRules.cronExp)
	for {
		select {
		case <- t.C:
			// todo  根据配置 写redis数据
			t.Reset(this.ContributionRules.cronExp)
		case <- this.Ctx.Done():
			return
		}
	}
}

func (this *Contribution) ChangeFids(topicIds []string) {
	this.topicIds = topicIds
}

func (this *Contribution) Stop() {
	this.cancel()
}


func (this *Contribution) GetName() (string) {
	return this.name
}

func (this *Contribution) reStart() {
	this.Stop()
	this.Ctx, this.cancel = context.WithCancel(context.Background())
	this.Start()
}

