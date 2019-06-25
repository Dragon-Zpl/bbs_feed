package kernel

import "time"

/*
	举报
*/
type reportThreadConf struct {
	reportCount int
}

type reportUserConf struct {
	reportCount int
}

type ThreadReportCheckEr struct {
	feedService *feedService
	reConf reportThreadConf
}

var ThreadReportCheck ThreadReportCheckEr

func (this *ThreadReportCheckEr) CheckThreadReport() {
	t := time.NewTimer(5 * time.Minute)
	for {
		select {
		case <-t.C:
			tids := this.GetReportTids()
			if len(tids) > 0 {
				this.feedService.mu.Lock()
				for _, agent := range this.feedService.Agents {
					if _, ok := agent.GetThis().(ThreadReport); ok {
						agent.(ThreadReport).AcceptSign(tids)
					}
				}
				this.feedService.mu.Unlock()
			}

			t.Reset(5 * time.Minute)
		}

	}
}


func (this *ThreadReportCheckEr) GetReportTids()(tids []string) {
	// todo 从配置 读出举报帖
	return
}

func (this *ThreadReportCheckEr) ChangeConf(conf reportThreadConf) {
	this.reConf = conf
}



type UserReportCheckEr struct {
	feedService *feedService
	reConf reportUserConf
}

var UserReportCheck UserReportCheckEr

func (this *UserReportCheckEr) CheckThreadReport() {
	t := time.NewTimer(5 * time.Minute)
	for {
		select {
		case <-t.C:
			tids := this.GetReportTids()
			if len(tids) > 0 {
				this.feedService.mu.Lock()
				for _, agent := range this.feedService.Agents {
					if _, ok := agent.GetThis().(ThreadReport); ok {
						agent.(ThreadReport).AcceptSign(tids)
					}
				}
				this.feedService.mu.Unlock()
			}

			t.Reset(5 * time.Minute)
		}

	}
}


func (this *UserReportCheckEr) GetReportTids()(tids []string) {
	// todo 从配置 读出举报帖
	return
}

func (this *UserReportCheckEr) ChangeConf(conf reportUserConf) {
	this.reConf = conf
}





