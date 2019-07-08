package kernel

import (
	"strings"
	"sync"
)

// 帖子举报
type ThreadReport interface {
	RemoveReportThread()
	AcceptSign([]string)
}

// 用户举报
type UserReport interface {
	RemoveReportUser()
	AcceptSign([]string)
}



type Agent interface {
	Init()
	Start()
	Stop()
	ChangeConf(interface{})
	GetName()string
	GetThis() interface{}
	ChangeFids([]string)
}

var FeedService *feedService


type feedService struct {
	Agents map[string]Agent
	mu *sync.Mutex
}

func NewFeedService(agents ...Agent) *feedService {
	agentsMap := make(map[string]Agent)
	for i := 0; i < len(agents); i ++ {
		agentsMap[agents[i].GetName()] = agents[i]
	}
	return &feedService{
		Agents: agentsMap,
		mu:     new(sync.Mutex),
	}
}

func (this *feedService) RegisterService(agent Agent){
	this.mu.Lock()
	defer this.mu.Unlock()
	agent.Init()
	agent.Start()
	this.Agents[agent.GetName()] = agent
	return
}

// 构造所有的agent
func (this *feedService) InitService() {
	for _, agent := range this.Agents {
		agent.Init()
	}
}

// 启动所有的agent
func (this *feedService) StartService() {
	for _, agent := range this.Agents {
		go agent.Start()
	}
}

func (this *feedService) StopService() {
	for _, agent := range this.Agents{
		agent.Stop()
	}
}

// 停止指定的agent
func (this *feedService) StopPusher(keys ...string) {
	for _, key := range keys {
		if _, ok := this.Agents[key]; ok {
			this.Agents[key].Stop()
		}
		delete(this.Agents, key)
	}
}


func (this *feedService) ChangeConf(conf interface{}) {
	for _, agent := range this.Agents {
		agent.ChangeConf(conf)
	}
}


func (this *feedService) ChangeFids(topicId string, topicIds []string) {
	for _, agent := range this.Agents {
		if strings.Split(agent.GetName(), "-")[0] == topicId {
			agent.ChangeFids(topicIds)
		}
	}
}


func NewSerivce() {
	FeedService = NewFeedService(CreateAgents()...)
	FeedService.InitService()
	FeedService.StartService()
}










