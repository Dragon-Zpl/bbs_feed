package contract

import (
	"bbs_feed/service"
	"bbs_feed/service/kernel/creater"
	"errors"
	"fmt"
	"strings"
	"sync"
)

type Agent interface {
	Init()
	Start()
	Stop()
	ChangeConf(string) error // 修改配置
	GetName() string
	GetThis() interface{}
	ChangeFids([]string)                              // 修改数据源
	AddTrait(id string, trait service.CallBlockTrait) // 添加额外信息
	Remover([]int)                                    // 删除agents redis 数据
}

//var FeedService *feedService

type FeedService struct {
	Agents map[string]Agent
	Mu     *sync.Mutex
}

func (this *FeedService) RegisterService(agents ...Agent) {
	this.Mu.Lock()
	defer this.Mu.Unlock()
	for _, agent := range agents {
		if _, ok := this.Agents[agent.GetName()]; ok {
			this.stopAgents(agent.GetName())
		}
		agent.Init()
		agent.Start()
		this.Agents[agent.GetName()] = agent
	}
	return
}

// 构造所有的agent
func (this *FeedService) InitService() {
	for _, agent := range this.Agents {
		agent.Init()
	}
}

// 启动所有的agent
func (this *FeedService) StartService() {
	for _, agent := range this.Agents {
		go agent.Start()
	}
}

func (this *FeedService) StopService() {
	for _, agent := range this.Agents {
		agent.Stop()
	}
}

// 停止指定的agent
func (this *FeedService) stopAgents(keys ...string) {
	for _, key := range keys {
		if _, ok := this.Agents[key]; ok {
			this.Agents[key].Stop()
		}
		delete(this.Agents, key)
	}
}

func (this *FeedService) RemovePusher(topicId string) {
	for _, agent := range this.Agents {
		if strings.Split(agent.GetName(), service.Separator)[0] == topicId {
			this.stopAgents(agent.GetName())
		}
	}
}

func (this *FeedService) ChangeConf(typ string, conf string) error {
	for _, agent := range this.Agents {
		if strings.Split(agent.GetName(), service.Separator)[1] == typ {
			if err := agent.ChangeConf(conf); err != nil {
				return err
			}
		}
	}
	return nil
}

func (this *FeedService) ChangeFids(topicId string, topicIds []string) {
	for _, agent := range this.Agents {
		if strings.Split(agent.GetName(), service.Separator)[0] == topicId {
			agent.ChangeFids(topicIds)
		}
	}
}

func (this *FeedService) Remove(agentName string, ids []int) error {
	if agent, ok := this.Agents[agentName]; ok {
		agent.Remover(ids)
		return nil
	} else {
		return errors.New(fmt.Sprintf("%s is not exist", agentName))
	}
}

var feedService *FeedService

func NewFeedService(agents ...Agent) *FeedService {
	agentsMap := make(map[string]Agent)
	for i := 0; i < len(agents); i++ {
		agentsMap[agents[i].GetName()] = agents[i]
	}
	return &FeedService{
		Agents: agentsMap,
		Mu:     new(sync.Mutex),
	}
}

func InstanceFeedService() *FeedService {
	return feedService
}

func InitFeedService() {
	feedService = NewFeedService(creater.InitGenAgents()...)
	feedService.InitService()
	feedService.StartService()
}
