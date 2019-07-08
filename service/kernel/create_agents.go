package kernel

import (
	"bbs_feed/model/feed_permission"
	"bbs_feed/service/kernel/call_block"
	"strings"
)

func CreateAgents()[]Agent{
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
	agents := make([]Agent, 0)
	for _, topic := range topics {
		topicIds := strings.Split(topic.TopicIds, ",")
		//agents = append(agents, call_block.NewContribution(topic.TopicId, topicIds),
		//						call_block.NewEssence(topic.TopicId, topicIds),
		//						call_block.NewHot(topic.TopicId, topicIds),
		//						call_block.NewWeekPopularity(topic.TopicId, topicIds))
		agents = append(agents,
			call_block.NewHot(topic.TopicId, topicIds))

	}
	return agents
}
