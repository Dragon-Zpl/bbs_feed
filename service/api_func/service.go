package api_func

import (
	"bbs_feed/boot"
	"bbs_feed/lib/stringi"
	"bbs_feed/model/feed_permission"
	"bbs_feed/service/kernel/contract"
	"bbs_feed/service/kernel/creater"
	"bbs_feed/v1/forms"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"strings"
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
func AddTopicService(form forms.TopicForm) error {
	if err1 := feed_permission.Insert(feed_permission.Model{
		TopicId:           stringi.ToInt(form.TopicId),
		HotThread:         stringi.ToInt(form.HotThread),
		Essence:           stringi.ToInt(form.Essence),
		TodayIntroduction: stringi.ToInt(form.TodayIntroduction),
		WeekPopularity:    stringi.ToInt(form.WeekPopularity),
		WeekContribution:  stringi.ToInt(form.WeekContribution),
		TopicIds:          form.TopicIds,
		IsUse:             1,
	}); err1 != nil {
		return err1
	}
	if agents, err2 := creater.GenAgents(form.TopicId); err2 != nil {
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

func GetRedisBlockDataService(topicId string, block string) (res_data []map[string]map[string]interface{}, err error){
	if err := feed_permission.GetBlock(topicId, block) ; err==nil{
		redis_key := "call_block_" + block + "_" + topicId
		data, err := boot.InstanceRedisCli(boot.CACHE).ZRange(redis_key,0,-1).Result()
		if err != nil{
			return nil, err
		}
		datas := "["+ strings.Join(data,",") + "]"
		res_data := make([]map[string]map[string]interface{}, 0, len(data))
		err = jsoniter.UnmarshalFromString(datas, &res_data)
		return res_data, nil
	} else {
		return nil, err
	}
}