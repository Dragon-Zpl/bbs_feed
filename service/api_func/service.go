package api_func

import (
	"bbs_feed/lib/helper"
	"bbs_feed/model/feed_conf"
	"bbs_feed/model/feed_permission"
	"bbs_feed/model/forum_thread"
	"bbs_feed/model/topic"
	"bbs_feed/model/topic_fid_relation"
	"bbs_feed/service"
	"bbs_feed/service/kernel/call_block"
	"bbs_feed/service/kernel/contract"
	"bbs_feed/service/kernel/creater"
	"bbs_feed/service/redis_ops"
	"bbs_feed/v1/forms"
	"errors"
	"fmt"
	"github.com/json-iterator/go"
	"log"
	"reflect"
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
	if err := feed_permission.UpdateFeedType(topicId, helper.Camel2Underline(feedTyp), isUse); err != nil {
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
	if err := forum_thread.UpdateDisplayorder(tids); err != nil {
		return err
	}
	creater.ThreadReportCheck.AcceptReportTids(tids)
	return nil
}

// 用户举报
func UserReportService(uids []int) {
	creater.UserReportCheck.AcceptReportUids(uids)
}

//删除指定调用块下的帖子或用户
func DelTopicDataService(agentName string, ids []int) error {
	return creater.InstanceFeedService().Remove(agentName, ids)
}

//添加额外信息
func AddCallBlockTraitService(form forms.TraitFrom) error {
	traitKey := fmt.Sprintf("call_block_%s_trait_%d", helper.Camel2Underline(form.FeedType), form.TopicId)
	if traitStr, err := jsoniter.MarshalToString(form.Trait); err != nil {
		return err
	} else {
		redis_ops.HSet(traitKey, strconv.Itoa(form.Id), traitStr, -1)
		//重启agent
		m, _ := feed_permission.GetOne(strconv.Itoa(form.TopicId))
		creater.InstanceFeedService().StopAgents(fmt.Sprintf("%d%s%s", form.TopicId, service.Separator, form.FeedType))
		creater.InstanceFeedService().RegisterService(creater.GenAgent(form.TopicId, form.FeedType, strings.Split(m.TopicIds, ",")))
	}
	return nil
}

// 获取结构体的字段
func GetFieldName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Tag.Get("json"))
	}
	return result
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

	block_datas["essenceRules"] = GetFieldName(essence)
	block_datas["hotRules"] = GetFieldName(hot)
	block_datas["newHotRules"] = GetFieldName(newHot)
	block_datas["introRules"] = GetFieldName(todayIntro)
	block_datas["weekPopularityRule"] = GetFieldName(weekPopularity)
	block_datas["contributionRules"] = GetFieldName(contribution)
	return block_datas
}

func GetTopicSerive() []map[string]interface{} {
	topicDatas := topic.GetAll()
	res_datas := make([]map[string]interface{}, 0, len(topicDatas))
	for _, data := range topicDatas {
		titles := make(map[string]interface{})
		titles["tid"] = data.Id
		titles["name"] = data.Title
		res_datas = append(res_datas, titles)
	}
	return res_datas
}
