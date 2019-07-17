package contract

import (
	"bbs_feed/model/report_num"
	"time"
)

const (
	CheckTime = 5 * time.Minute
)

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
	ReportChan chan []int
}

func (this *ThreadRep) RemoveReportThread(f func([]int)) {
	for {
		select {
		case tids := <-this.ReportChan:
			f(tids)
		}
	}
}

func (this *ThreadRep) AcceptSign(tids []int) {
	this.ReportChan <- tids
}

func CreateThreadReport() ThreadReport {
	return &ThreadRep{
		ReportChan: make(chan []int, 10),
	}
}

type ReportThreadConf struct {
	ReportCount int `form:"reportCount" binding:"required"`
}

type ReportUserConf struct {
	ReportCount int `form:"reportCount" binding:"required"`
}

type ThreadReportCheckEr struct {
	FeedService *FeedService
	ReConf      ReportThreadConf
	ReportTids  chan []int
}

func (this *ThreadReportCheckEr) CheckThreadReport() {
	t := time.NewTimer(CheckTime)
	for {
		select {
		case <-t.C:
			tids := this.GetReportTids()
			if len(tids) > 0 {
				this.seedReportTids(tids)
			}
			t.Reset(CheckTime)
		case tids := <-this.ReportTids:
			this.seedReportTids(tids)
		}
	}
}

func (this *ThreadReportCheckEr) seedReportTids(tids []int) {
	this.FeedService.Mu.Lock()
	for _, agent := range this.FeedService.Agents {
		if _, ok := agent.GetThis().(ThreadReport); ok {
			agent.(ThreadReport).AcceptSign(tids)
		}
	}
	this.FeedService.Mu.Unlock()
}

// 处理举报贴接口
func (this *ThreadReportCheckEr) AcceptReportTids(tids []int) {
	this.ReportTids <- tids
}

func (this *ThreadReportCheckEr) GetReportTids() (tids []int) {
	results := report_num.GetAll(CheckTime)
	for _, result := range results {
		if result.Num > this.ReConf.ReportCount {
			tids = append(tids, result.ThreadId)
		}
	}
	return
}

func (this *ThreadReportCheckEr) ChangeConf(conf ReportThreadConf) {
	this.ReConf = conf
}

type UserReportCheckEr struct {
	FeedService *FeedService
	ReConf      ReportUserConf
	ReportUids  chan []int
}

func (this *UserReportCheckEr) CheckUserReport() {
	t := time.NewTimer(CheckTime)
	for {
		select {
		case <-t.C:
			uids := this.GetReportUids()
			if len(uids) > 0 {
				this.seedReportUids(uids)
			}

			t.Reset(CheckTime)
		case uids := <-this.ReportUids:
			this.seedReportUids(uids)
		}

	}
}

func (this *UserReportCheckEr) seedReportUids(tids []int) {
	this.FeedService.Mu.Lock()
	for _, agent := range this.FeedService.Agents {
		if _, ok := agent.GetThis().(UserReport); ok {
			agent.(UserReport).AcceptSign(tids)
		}
	}
	this.FeedService.Mu.Unlock()
}

func (this *UserReportCheckEr) GetReportUids() (tids []int) {
	// todo 从配置 读出举报帖
	return
}

func (this *UserReportCheckEr) ChangeConf(conf ReportUserConf) {
	this.ReConf = conf
}

// 处理举报贴接口
func (this *UserReportCheckEr) AcceptReportUids(uids []int) {
	this.ReportUids <- uids
}
