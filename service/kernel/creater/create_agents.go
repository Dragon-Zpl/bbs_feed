package creater

import (
	"bbs_feed/model/feed_permission"
	"bbs_feed/service/kernel/call_block"
	"bbs_feed/service/kernel/contract"
	"reflect"
	"strings"
)

/*
	topicId == 0 即为全局
*/

type AgentGen func(int, []string) contract.Agent

var AgentMapping = map[string]AgentGen {
	"hot_thread": HotThread(),
	"essence": Essence(),
	"new_hot": NewHot(),
	"today_introduction": TodayIntroduction(),
	"week_popularity": WeekContribution(),
	"week_contribution": WeekContribution(),
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
		return nil
	}
}

// 最新最热
func NewHot() AgentGen {
	return func(topicId int, topicIds []string) contract.Agent {
		return nil
	}
}

// 今日导读
func TodayIntroduction() AgentGen {
	return func(topicId int, topicIds []string) contract.Agent {
		return nil
	}
}


func WeekPopularity() AgentGen {
	return func(topicId int, topicIds []string) contract.Agent {
		return nil
	}
}


func WeekContribution() AgentGen {
	return func(topicId int, topicIds []string) contract.Agent {
		return nil
	}
}
func CreateAgents()[]contract.Agent{
	topics := feed_permission.GetAll()
	agents := make([]contract.Agent, 0)

	for _, one := range topics {
		topicId := one.TopicId
		topicIds := strings.Split(one.TopicIds, ",")

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

	//topics = []*feed_permission.Model{
	//	&feed_permission.Model{
	//		TopicId:        2,
	//		HotThread:         "",
	//		Essence:        "",
	//		WeekPopularity: "",
	//		Contribution:   "",
	//		TopicIds:       "2, 3, 4",
	//		IsUse:          1,
	//	},
	//}
	//for _, topic := range topics {
	//	topicIds := strings.Split(topic.TopicIds, ",")
	//	agents = append(agents, AllTypeAgents(topic.TopicId, topicIds)...)
	//
	//}
	return agents
}

func CreateAgent(topicId string) ([]contract.Agent, error){
	topic, err := feed_permission.GetOne(topicId)
	if err != nil {
		return nil, err
	}
	agents := make([]contract.Agent, 0)
	tid := topic.TopicId
	topicIds := strings.Split(topic.TopicIds, ",")

	t := reflect.TypeOf(topic).Elem()
	v := reflect.ValueOf(topic).Elem()
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

