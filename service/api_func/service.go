package api_func

import (
	"bbs_feed/model/feed_conf"
	"bbs_feed/model/feed_permission"
	"bbs_feed/model/topic"
	"bbs_feed/model/topic_fid_relation"
	"bbs_feed/service"
	"bbs_feed/service/kernel/call_block"
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
	fid := topic_fid_relation.GetFids([]string{strconv.Itoa(form.TopicId)})
	if err1 := feed_permission.Insert(feed_permission.Model{
		TopicId:           form.TopicId,
		Fid:               fid[0],
		Hot:               form.Hot,
		NewHot:            form.NewHot,
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

		if agent := creater.GenAgent(topicId, feedTyp, topicIds); agent != nil {
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
			//TODO 重启服务?
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
func ThreadReportService(tids []int) error {
	//if err := forum_thread.UpdateDisplayorder(tids); err != nil {
	//	return err
	//}
	creater.ThreadReportCheck.AcceptReportTids(tids)
	return nil
}

// 用户举报
func UserReportService(uids []int) {
	creater.UserReportCheck.AcceptReportUids(uids)
}

// 获取板块可改字段
func GetFeedConfUseSerive() map[string]interface{} {
	block_datas := make(map[string]interface{})
	var (
		hot            call_block.HotRules
		essence        call_block.EssenceRules
		contribution   call_block.ContributionRules
		weekPopularity call_block.WeekPopularityRule
		newHot         call_block.NewHotRules
		todayIntro     call_block.IntroRules
	)

	block_datas["essenceRules"] = essence
	block_datas["hotRules"] = hot
	block_datas["newHotRules"] = newHot
	block_datas["introRules"] = todayIntro
	block_datas["weekPopularityRule"] = weekPopularity
	block_datas["contributionRules"] = contribution
	return block_datas
}

func GetTopicSerive() map[string]map[string]interface{} {
	topicDatas := topic.GetAll()
	topicDatasMap := make(map[string]string)
	for _,data := range topicDatas{
		topicDatasMap[strconv.Itoa(data.Id)] = data.Title
	}
	preMisDatas := feed_permission.GetAll()
	preMissDataMap := make(map[string]*feed_permission.Model)
	for _,data := range preMisDatas{
		preMissDataMap[strconv.Itoa(data.TopicId)] = data
	}
	res_datas := make(map[string]map[string]interface{})
	for _,data := range topicDatas{
		if _,ok := preMissDataMap[strconv.Itoa(data.Id)]; !ok{
			titles := make(map[string]interface{})
			titles["titles"] = ""
			titles["isuse"] = 0
			res_datas[data.Title] = titles
			continue
		}
		preMisData := preMissDataMap[strconv.Itoa(data.Id)]
		topicids := strings.Split(preMisData.TopicIds,",")
		all_titles := make([]string,0)
		for _,topicid := range topicids{
			if data,ok :=topicDatasMap[topicid]; ok{
				all_titles = append(all_titles, data)
			}
		}
		titles := make(map[string]interface{})
		titles["titles"] = all_titles
		titles["isuse"] = preMisData.IsUse
		res_datas[data.Title] = titles
	}
	return res_datas
}
