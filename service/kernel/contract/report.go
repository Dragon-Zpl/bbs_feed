package contract

import (
	"time"
)

/*
	举报
*/

// 帖子举报
type ThreadReport interface {
	RemoveReportThread(func([]int))
	AcceptSign([]int)
}

// 用户举报
type UserReport interface {
	RemoveReportUser()
	AcceptSign([]int)
}

type ThreadRep struct {
	reportChan chan []int
}

func (this *ThreadRep) RemoveReportThread(f func([]int)) {
	for {
		select {
		case tids := <-this.reportChan:
			f(tids)
		}
	}
}

func (this *ThreadRep) AcceptSign(tids []int) {
	this.reportChan <- tids
}

func CreateThreadReport() ThreadReport {
	return &ThreadRep{
		reportChan: make(chan []int, 10),
	}
}

type ReportThreadConf struct {
	ReportCount int `form:"reportCount" binding:"required"`
}

type ReportUserConf struct {
	ReportCount int `form:"reportCount" binding:"required"`
}

type ThreadReportCheckEr struct {
	feedService *FeedService
	reConf      ReportThreadConf
	reportTids  chan []int
}

var ThreadReportCheck *ThreadReportCheckEr

func NewThreadReportCheckEr() {
	ThreadReportCheck = &ThreadReportCheckEr{
		feedService: feedService,
		reConf:      ReportThreadConf{},
		reportTids:  make(chan []int, 10),
	}
	go ThreadReportCheck.CheckThreadReport()
}

func (this *ThreadReportCheckEr) CheckThreadReport() {
	t := time.NewTimer(5 * time.Minute)
	for {
		select {
		case <-t.C:
			tids := this.GetReportTids()
			if len(tids) > 0 {
				this.seedReportTids(tids)
			}

			t.Reset(5 * time.Minute)
		case tids := <-this.reportTids:
			this.seedReportTids(tids)
		}
	}
}

func (this *ThreadReportCheckEr) seedReportTids(tids []int) {
	this.feedService.Mu.Lock()
	for _, agent := range this.feedService.Agents {
		if _, ok := agent.GetThis().(ThreadReport); ok {
			agent.(ThreadReport).AcceptSign(tids)
		}
	}
	this.feedService.Mu.Unlock()
}

// 处理举报贴接口
func (this *ThreadReportCheckEr) AcceptReportTids(tids []int) {
	this.reportTids <- tids
}

func (this *ThreadReportCheckEr) GetReportTids() (tids []int) {
	// todo 从配置 读出举报帖
	return
}

func (this *ThreadReportCheckEr) ChangeConf(conf ReportThreadConf) {
	this.reConf = conf
}

type UserReportCheckEr struct {
	feedService *FeedService
	reConf      ReportUserConf
	reportUids  chan []int
}

var UserReportCheck *UserReportCheckEr

func NewUserReportCheck() {
	UserReportCheck = &UserReportCheckEr{
		feedService: feedService,
		reConf:      ReportUserConf{},
		reportUids:  make(chan []int, 10),
	}
	go UserReportCheck.CheckUserReport()
}

func (this *UserReportCheckEr) CheckUserReport() {
	t := time.NewTimer(5 * time.Minute)
	for {
		select {
		case <-t.C:
			uids := this.GetReportUids()
			if len(uids) > 0 {
				this.seedReportUids(uids)
			}

			t.Reset(5 * time.Minute)
		case uids := <-this.reportUids:
			this.seedReportUids(uids)
		}

	}
}

func (this *UserReportCheckEr) seedReportUids(tids []int) {
	this.feedService.Mu.Lock()
	for _, agent := range this.feedService.Agents {
		if _, ok := agent.GetThis().(UserReport); ok {
			agent.(UserReport).AcceptSign(tids)
		}
	}
	this.feedService.Mu.Unlock()
}

func (this *UserReportCheckEr) GetReportUids() (tids []int) {
	// todo 从配置 读出举报帖
	return
}

func (this *UserReportCheckEr) ChangeConf(conf ReportUserConf) {
	this.reConf = conf
}

// 处理举报贴接口
func (this *UserReportCheckEr) AcceptReportUids(uids []int) {
	this.reportUids <- uids
}
