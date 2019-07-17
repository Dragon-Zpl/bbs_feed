package creater

import (
	"bbs_feed/service/kernel/contract"
	"sync"
)

//默认举报次数界限
const ReportCountLimit = 3

//feed service
var feedService *contract.FeedService

func NewFeedService(agents ...contract.Agent) *contract.FeedService {
	agentsMap := make(map[string]contract.Agent)
	for i := 0; i < len(agents); i++ {
		agentsMap[agents[i].GetName()] = agents[i]
	}
	return &contract.FeedService{
		Agents: agentsMap,
		Mu:     new(sync.Mutex),
	}
}

func InstanceFeedService() *contract.FeedService {
	return feedService
}

func InitFeedService() {
	feedService = NewFeedService(InitGenAgents()...)
	feedService.InitService()
	feedService.StartService()
}

//帖子举报service
var ThreadReportCheck *contract.ThreadReportCheckEr

func NewThreadReportCheckEr() {
	ThreadReportCheck = &contract.ThreadReportCheckEr{
		FeedService: feedService,
		ReConf:      contract.ReportThreadConf{ReportCount: ReportCountLimit},
		ReportTids:  make(chan []int, 10),
	}
	go ThreadReportCheck.CheckThreadReport()
}

//用户举报service
var UserReportCheck *contract.UserReportCheckEr

func NewUserReportCheckEr() {
	UserReportCheck = &contract.UserReportCheckEr{
		FeedService: feedService,
		ReConf:      contract.ReportUserConf{},
		ReportUids:  make(chan []int, 10),
	}
	go UserReportCheck.CheckUserReport()
}

func InitService() {
	InitFeedService()
	NewThreadReportCheckEr()
	NewUserReportCheckEr()
}
