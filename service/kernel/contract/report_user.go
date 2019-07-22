package contract

import (
	"bbs_feed/lib/helper"
	"bbs_feed/model/common_member_crime"
	"time"
)

/*处理违规用户*/

type UserReport interface {
	RemoveReportUser(func([]int))
	AcceptSign([]int)
}

type ReportUserConf struct {
	ReportCount int `form:"reportCount"`
}

type UserRep struct {
	ReportChan chan []int
}

func (this *UserRep) RemoveReportUser(f func([]int)) {
	for {
		select {
		case uids := <-this.ReportChan:
			f(uids)
		}
	}
}

func (this *UserRep) AcceptSign(uids []int) {
	this.ReportChan <- uids
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

func (this *UserReportCheckEr) seedReportUids(uids []int) {
	this.FeedService.Mu.Lock()
	for _, agent := range this.FeedService.Agents {
		if _, ok := agent.GetThis().(UserReport); ok {
			agent.(UserReport).AcceptSign(uids)
		}
	}
	this.FeedService.Mu.Unlock()
}

func (this *UserReportCheckEr) GetReportUids() (uids []int) {
	results := common_member_crime.GetAll(CheckTime)
	allUids := make([]int, 0, len(results))
	for _, result := range results {
		allUids = append(allUids, result.Uid)
	}
	uidsMap := helper.SameElementCount(allUids)
	for uid, count := range uidsMap {
		if count > this.ReConf.ReportCount {
			uids = append(uids, uid)
		}
	}
	return
}

func (this *UserReportCheckEr) ChangeConf(conf ReportUserConf) {
	this.ReConf = conf
}

// 接收违规用户uid
func (this *UserReportCheckEr) AcceptReportUids(uids []int) {
	this.ReportUids <- uids
}
