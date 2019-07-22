package contract

import (
	"bbs_feed/model/report_num"
	"time"
)

/*举报帖子*/

const (
	CheckTime = 5 * time.Minute
)

type ThreadReport interface {
	RemoveReportThread(func([]int))
	AcceptSign([]int)
}

type ReportThreadConf struct {
	ReportCount int `form:"reportCount"`
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

// 接收举报的帖子
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
