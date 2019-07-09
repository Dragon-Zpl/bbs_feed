package creater

import (
	"bbs_feed/model/feed_permission"
	"bbs_feed/service/kernel/call_block"
	"bbs_feed/service/kernel/contract"

	"strings"
)

func CreateAgents()[]contract.Agent{
	topics := feed_permission.GetAll()
	topics = []*feed_permission.Model{
		&feed_permission.Model{
			TopicId:        2,
			Hot:            "",
			Essence:        "",
			WeekPopularity: "",
			Contribution:   "",
			TopicIds:       "2, 3, 4",
			IsUse:          1,
		},
	}
	agents := make([]contract.Agent, 0)
	for _, topic := range topics {
		topicIds := strings.Split(topic.TopicIds, ",")
		agents = append(agents, AllTypeAgents(topic.TopicId, topicIds)...)

	}
	return agents
}

func CreateAgent(topicId string) ([]contract.Agent, error){
	topic, err := feed_permission.GetOne(topicId)
	if err != nil {
		return nil, err
	}
	topicIds := strings.Split(topic.TopicIds, ",")
	agents := AllTypeAgents(topic.TopicId, topicIds)
	return agents, nil
}


func AllTypeAgents(topicId int, topicIds []string) []contract.Agent {
	agents := make([]contract.Agent, 0)
	//agents = append(agents, call_block.NewContribution(topic.TopicId, topicIds),
	//						call_block.NewEssence(topic.TopicId, topicIds),
	//						call_block.NewHot(topic.TopicId, topicIds),
	//						call_block.NewWeekPopularity(topic.TopicId, topicIds))
	agents = append(agents,
		call_block.NewHot(topicId, topicIds, contract.CreateThreadReport()))

	return agents
}