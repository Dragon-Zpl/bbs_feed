package creater

import (
	"bbs_feed/model/feed_permission"
	"bbs_feed/service/data_source"
	"bbs_feed/service/kernel/call_block"
	"bbs_feed/service/kernel/contract"
	"reflect"
	"strings"
	"sync"
)

/*
	topicId == 0 即为全局
*/

type AgentGen func(int, []string) contract.Agent

var once sync.Once

var AgentMapping = map[string]AgentGen{
	"hot":               HotThread(),
	"essence":           Essence(),
	"newHot":            NewHot(),
	"todayIntroduction": TodayIntroduction(),
	"weekPopularity":    WeekPopularity(),
	"weekContribution":  WeekContribution(),
}

// 热门贴
func HotThread() AgentGen {
	return func(topicId int, topicIds []string) contract.Agent {
		return call_block.NewHot(topicId, topicIds)
	}
}

// 精华贴
func Essence() AgentGen {
	return func(topicId int, topicIds []string) contract.Agent {
		return call_block.NewEssence(topicId, topicIds)
	}
}

// 最新最热
func NewHot() AgentGen {
	return func(topicId int, topicIds []string) contract.Agent {
		return call_block.NewNewHots(topicId, topicIds)
	}
}

// 今日导读
func TodayIntroduction() AgentGen {
	return func(topicId int, topicIds []string) contract.Agent {
		return call_block.NewTodayIntro(topicId, topicIds)
	}
}

//本周人气榜
func WeekPopularity() AgentGen {
	return func(topicId int, topicIds []string) contract.Agent {
		return nil
	}
}

//本周贡献榜
func WeekContribution() AgentGen {
	return func(topicId int, topicIds []string) contract.Agent {
		return nil
	}
}

// agents 的生成器 用于启动时
func InitGenAgents() []contract.Agent {
	once.Do(call_block.InitConfs)
	topics := feed_permission.GetAll()
	agents := make([]contract.Agent, 0)

	for _, one := range topics {
		topicId := one.TopicId
		topicIds := strings.Split(one.TopicIds, ",")
		topicIds = data_source.CheckAndGetCurTopics(topicIds)
		t := reflect.TypeOf(one).Elem()
		v := reflect.ValueOf(one).Elem()
		for i := 0; i < t.NumField(); i++ {
			tag := t.Field(i).Tag.Get("json")
			if _, ok := AgentMapping[tag]; ok {
				if v.Field(i).Int() == 1 {
					agents = append(agents, AgentMapping[tag](topicId, topicIds))
				}
			}
		}
	}

	return agents
}

// agents 的生成器 用于topic
func GenAgents(topicId string) ([]contract.Agent, error) {
	once.Do(call_block.InitConfs)
	topic, err := feed_permission.GetOne(topicId)
	if err != nil {
		return nil, err
	}
	agents := make([]contract.Agent, 0)
	tid := topic.TopicId
	topicIds := strings.Split(topic.TopicIds, ",")
	topicIds = data_source.CheckAndGetCurTopics(topicIds)

	t := reflect.TypeOf(&topic).Elem()
	v := reflect.ValueOf(&topic).Elem()
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")
		if _, ok := AgentMapping[tag]; ok {
			if v.Field(i).Int() == 1 {
				agents = append(agents, AgentMapping[tag](tid, topicIds))
			}
		}
	}
	return agents, nil
}

func GenAgent(topicId int, feedTyp string, topicIds []string) contract.Agent {
	topicIds = data_source.CheckAndGetCurTopics(topicIds)

	if f, ok := AgentMapping[feedTyp]; ok {
		return f(topicId, topicIds)
	}
	return nil
}
