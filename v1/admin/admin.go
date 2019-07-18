package admin

import (
	"bbs_feed/lib/feed_errors"
	"bbs_feed/lib/helper"
	"bbs_feed/service"
	"bbs_feed/service/api_func"
	"bbs_feed/service/kernel/contract"
	"bbs_feed/service/redis_ops"
	"bbs_feed/v1/forms"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"strings"
	"time"
)

func Mapping(prefix string, app *gin.Engine) {
	admin := app.Group(prefix)
	admin.POST("/topic", feed_errors.MdError(AddTopic))
	admin.PUT("/topic", feed_errors.MdError(UpdateTopic))
	admin.PUT("/agent", feed_errors.MdError(UpdateAgent))
	admin.PUT("/topic/topic_ids", feed_errors.MdError(UpdateTopicIds))

	admin.POST("/feed/conf", feed_errors.MdError(AddFeedTypeConf))
	admin.PUT("/feed/conf", feed_errors.MdError(UpdateFeedTypeConf))

	admin.PUT("/thread/report_conf", feed_errors.MdError(UpdateThreadReportConf))
	admin.PUT("/user/report_conf", feed_errors.MdError(UpdateUserReportConf))
	admin.DELETE("/call_back", feed_errors.MdError(DelTopicData))

	admin.POST("/report/thread", feed_errors.MdError(ThreadReport))
	admin.POST("/report/user", feed_errors.MdError(UserReport))
	admin.POST("/trait", feed_errors.MdError(AddCallBlockTrait))
}

// 添加topic
func AddTopic(ctx *gin.Context) error {
	var topicForm forms.TopicForm
	if err := ctx.ShouldBind(&topicForm); err != nil {
		return errors.New("params_error")
	}
	if err := api_func.AddTopicService(topicForm); err != nil {
		return err
	}
	ctx.JSON(helper.Success())
	return nil
}

// 启用、关闭topic
func UpdateTopic(ctx *gin.Context) error {
	var updateTopicForm forms.UpdateTopicForm
	if err := ctx.ShouldBind(&updateTopicForm); err != nil {
		return errors.New("params_error")
	}
	if err := api_func.UpdateTopicService(updateTopicForm.TopicId, updateTopicForm.IsUse); err != nil {
		return err
	}
	ctx.JSON(helper.Success())
	return nil
}

//启用、关闭Agent
func UpdateAgent(ctx *gin.Context) error {
	var agentForm forms.AgentForm
	if err := ctx.ShouldBind(&agentForm); err != nil {
		return errors.New("params_error")
	}
	if err := api_func.UpdateAgentService(agentForm.TopicId, agentForm.FeedType, agentForm.IsUse); err != nil {
		return err
	} else {
		ctx.JSON(helper.Success())
		return nil
	}
	return nil
}

// 修改topic数据源topicIds
func UpdateTopicIds(ctx *gin.Context) error {
	var topicDataSourceForm forms.TopicDataSourceForm
	if err := ctx.ShouldBind(&topicDataSourceForm); err != nil {
		return errors.New("params_error")
	}
	api_func.UpdateTopicIdsService(topicDataSourceForm.TopicId, topicDataSourceForm.TopicIds)
	ctx.JSON(helper.Success())
	return nil
}

//添加调用块配置
func AddFeedTypeConf(ctx *gin.Context) error {
	var feedTypeConfForm forms.FeedTypeConfForm
	if err := ctx.ShouldBind(&feedTypeConfForm); err != nil {
		return errors.New("params_error")
	}
	if err := api_func.AddFeedTypeConfService(feedTypeConfForm.FeedType, feedTypeConfForm.Conf); err != nil {
		return err
	}
	return nil
}

// 修改调用块配置
func UpdateFeedTypeConf(ctx *gin.Context) error {
	var feedTypeConfForm forms.FeedTypeConfForm
	if err := ctx.ShouldBind(&feedTypeConfForm); err != nil {
		return errors.New("params_error")
	}
	if err := api_func.UpdateFeedTypeConfService(feedTypeConfForm.FeedType, feedTypeConfForm.Conf); err != nil {
		return errors.New("conf_error")
	}
	ctx.JSON(helper.Success())
	return nil
}

// 修改帖子举报规则
func UpdateThreadReportConf(ctx *gin.Context) error {
	var reportThreadConf contract.ReportThreadConf
	if err := ctx.ShouldBind(&reportThreadConf); err != nil {
		return errors.New("params_error")
	}
	api_func.UpdateThreadReportConfService(reportThreadConf)
	ctx.JSON(helper.Success())
	return nil
}

// 修改用户举报规则
func UpdateUserReportConf(ctx *gin.Context) error {
	var reportUserConf contract.ReportUserConf
	if err := ctx.ShouldBind(&reportUserConf); err != nil {
		return errors.New("params_error")
	}
	api_func.UpdateUserReportConfService(reportUserConf)
	ctx.JSON(helper.Success())
	return nil
}

//删除指定板块下的数据(tid、uid)
func DelTopicData(ctx *gin.Context) error {
	var delTopicDataFrom forms.DelTopicDataFrom
	if err := ctx.ShouldBind(&delTopicDataFrom); err != nil {
		return errors.New("params_error")
	}
	agentName := fmt.Sprintf("%s%s%s", delTopicDataFrom.TopicId, service.Separator, delTopicDataFrom.FeedType)
	ids := strings.Split(delTopicDataFrom.Ids, ",")
	if err := api_func.DelTopicDataService(agentName, helper.ArrayStrToInt(ids)); err != nil {
		return err
	}
	ctx.JSON(helper.Success())
	return nil
}

// 帖子举报
func ThreadReport(ctx *gin.Context) error {
	var threadReportForm forms.ThreadReportForm
	if err := ctx.ShouldBind(&threadReportForm); err != nil {
		return errors.New("params_error")
	}
	tids := strings.Split(threadReportForm.ThreadIds, ",")
	api_func.ThreadReportService(helper.ArrayStrToInt(tids))
	ctx.JSON(helper.Success())
	return nil
}

// 用户举报
//TODO 待确定
func UserReport(ctx *gin.Context) error {
	var userReportForm forms.UserReportForm
	if err := ctx.ShouldBind(&userReportForm); err != nil {
		return errors.New("params_error")
	}
	uids := strings.Split(userReportForm.UserIds, ",")
	api_func.UserReportService(helper.ArrayStrToInt(uids))
	ctx.JSON(helper.Success())
	return nil
}

//添加额外信息
func AddCallBlockTrait(ctx *gin.Context) error {
	var traitFrom forms.TraitFrom
	if err := ctx.ShouldBind(&traitFrom); err != nil {
		return errors.New("params_error")
	}
	redisKey := fmt.Sprintf("call_block_%s_trait", traitFrom.FeedType)
	if !redis_ops.KeyExist(redisKey) {
		return errors.New("redis_key_notexist")
	}
	if traitStr, err := jsoniter.MarshalToString(traitFrom.Trait); err != nil {
		return err
	} else {
		redis_ops.HSet(redisKey, traitFrom.Id, traitStr, time.Duration(traitFrom.Exp)*time.Hour)
	}
	ctx.JSON(helper.Success())
	return nil
}
