package kernel

import (
	"bbs_feed/model/feed_permission"
	"strings"
)

func CreateAgents()[]Agent{
	topics := feed_permission.GetAll()
	agents := make([]Agent, 0)
	for _, topic := range topics {
		topicIds := strings.Split(topic.TopicIds, ",")
		agents = append(agents, NewContribution(topic.TopicId, topicIds))
		agents = append(agents, NewEssence(topic.TopicId, topicIds))
		agents = append(agents, NewHot(topic.TopicId, topicIds))
		agents = append(agents, NewWeekPopularity(topic.TopicId, topicIds))
	}
	return agents
}
