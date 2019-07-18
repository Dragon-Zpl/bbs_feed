package api_func

import (
	"bbs_feed/model/feed_conf"
	"bbs_feed/model/feed_permission"
	"bbs_feed/service"
	"bbs_feed/service/kernel/contract"
	"bbs_feed/service/kernel/creater"
	"bbs_feed/v1/forms"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	NotUse = 0
	IsUse  = 1
)

// 增加topic
func AddTopicService(form forms.TopicForm) error {
	if err1 := feed_permission.Insert(feed_permission.Model{
		TopicId:           form.TopicId,
		HotThread:         form.HotThread,
		Essence:           form.Essence,
		TodayIntroduction: form.TodayIntroduction,
		WeekPopularity:    form.WeekPopularity,
		WeekContribution:  form.WeekContribution,
		TopicIds:          form.TopicIds,
		IsUse:             1,
	}); err1 != nil {
		return err1
	}
	if agents, err2 := creater.GenAgents(strconv.Itoa(form.TopicId)); err2 != nil {
		return errors.New("topic_not_deploy")
	} else {
		creater.InstanceFeedService().RegisterService(agents...)
	}
	return nil
}

// 启用、关闭topic
func UpdateTopicService(topicId int, isUse int) error {
	if err1 := feed_permission.UpdateIsUse(topicId, isUse); err1 != nil {
		return err1
	}
	if isUse == NotUse { //关闭
		creater.InstanceFeedService().RemovePusher(strconv.Itoa(topicId))
	} else if isUse == IsUse { //启用
		if agents, err2 := creater.GenAgents(strconv.Itoa(topicId)); err2 != nil {
			return errors.New("topic_not_deploy")
		} else {
			creater.InstanceFeedService().RegisterService(agents...)
		}
	}
	return nil
}

//启用、关闭agent
func UpdateAgentService(topicId int, feedTyp string, isUse int) error {
	if err := feed_permission.UpdateFeedType(topicId, feedTyp, isUse); err != nil {
		return err
	}
	if isUse == IsUse { //启用
		m, _ := feed_permission.GetOne(strconv.Itoa(topicId))
		topicIds := strings.Split(m.TopicIds, ",")

		if agent := creater.GenAgent(topicId, feedTyp, topicIds); agent == nil {
			creater.InstanceFeedService().RegisterService(agent)
		} else {
			return errors.New("params_error")
		}
	} else if isUse == NotUse { //关闭
		creater.InstanceFeedService().StopAgents(fmt.Sprintf("%d%s%s", topicId, service.Separator, feedTyp))
	}
	return nil
}

// topic 数据源改变
func UpdateTopicIdsService(topicId string, topicIds string) error {
	if err := feed_permission.UpdateTopicIds(topicId, topicIds); err != nil {
		return err
	}
	creater.InstanceFeedService().ChangeFids(topicId, strings.Split(topicIds, ","))
	return nil
}

//添加调用块配置
func AddFeedTypeConfService(typ string, conf string) error {
	if _, err1 := feed_conf.GetOne(typ); err1 != nil {
		if err2 := feed_conf.Insert(feed_conf.Model{
			Name:  typ,
			Conf:  conf,
			IsUse: 1,
		}); err2 != nil {
			return err2
		} else {
			return creater.InstanceFeedService().ChangeConf(typ, conf)
		}
	} else {
		return errors.New("feed_type_exist")
	}
}

// 修改调用块配置
func UpdateFeedTypeConfService(typ string, conf string) error {
	if err := feed_conf.UpdateConf(typ, conf); err != nil {
		return err
	}
	return creater.InstanceFeedService().ChangeConf(typ, conf)
}

// 修改帖子举报规则
func UpdateThreadReportConfService(conf contract.ReportThreadConf) {
	creater.ThreadReportCheck.ChangeConf(conf)
}

// 修改用户举报规则
func UpdateUserReportConfService(conf contract.ReportUserConf) {
	creater.UserReportCheck.ChangeConf(conf)
}

// 帖子举报
func ThreadReportService(tids []int) {
	creater.ThreadReportCheck.AcceptReportTids(tids)
}

// 用户举报
func UserReportService(uids []int) {
	creater.UserReportCheck.AcceptReportUids(uids)
}

func DelTopicDataService(agentName string, ids []int) error {
	return creater.InstanceFeedService().Remove(agentName, ids)
}
