package api_func

import (
	"bbs_feed/model/feed_permission"
	"bbs_feed/service/kernel/contract"
	"bbs_feed/service/kernel/creater"
	"errors"
)

// 调用块配置改变
func FeedTypeConfChangeService(typ string, conf string) error {
	return creater.InstanceFeedService().ChangeConf(typ, conf)
}

// topic 数据源改变
func TopicDataSourceChangeService(topicId string, topicIds []string) {
	creater.InstanceFeedService().ChangeFids(topicId, topicIds)
}

// 增加topic
func AddTopicService(topicId string) error {
	if agents, err := creater.GenAgents(topicId); err != nil {
		return errors.New("topic_not_deploy")
	} else {
		creater.InstanceFeedService().RegisterService(agents...)
	}
	return nil
}

//添加agent
func AddAgentService(topicId int, feedTyp string, topicIds []string) error {
	if agent := creater.GenAgent(topicId, feedTyp, topicIds); agent == nil {
		creater.InstanceFeedService().RegisterService(agent)
		return nil
	} else {
		return errors.New("params_error")
	}
}

// 删除topic
func DelTopicService(topicId string) {
	feed_permission.UpdateIsUse(topicId, 0)
	creater.InstanceFeedService().RemovePusher(topicId)
}

// 修改帖子举报规则
func ChangeThreadReportConfService(conf contract.ReportThreadConf) {
	creater.ThreadReportCheck.ChangeConf(conf)
}

// 帖子举报
func ThreadReportService(tids []int) {
	creater.ThreadReportCheck.AcceptReportTids(tids)
}

// 修改用户举报规则
func ChangeUserReportConfService(conf contract.ReportUserConf) {
	creater.UserReportCheck.ChangeConf(conf)
}

// 用户举报
func UserReportService(uids []int) {
	creater.UserReportCheck.AcceptReportUids(uids)
}

func DelTopicDataService(agentName string, ids []int) error {
	return creater.InstanceFeedService().Remove(agentName, ids)
}
