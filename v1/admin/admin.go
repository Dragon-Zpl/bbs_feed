package admin

import (
	"bbs_feed/lib/feed_errors"
	"bbs_feed/lib/helper"
	"bbs_feed/service"
	"bbs_feed/service/api_func"
	"bbs_feed/service/kernel/contract"
	"bbs_feed/service/redis_ops"
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
	admin.POST("/agnet", feed_errors.MdError(AddAgent))
	admin.DELETE("/topic/:id", feed_errors.MdError(DelTopic))
	admin.POST("/topic/conf", feed_errors.MdError(FeedTypeConfChange))
	admin.POST("/topic/fids", feed_errors.MdError(TopicDataSourceChange))
	admin.POST("/report/thread/conf", feed_errors.MdError(ChangeThreadReportConf))
	admin.POST("/report/user/conf", feed_errors.MdError(ChangeUserReportConf))
	admin.DELETE("/call_back", feed_errors.MdError(DelTopicData))
	admin.POST("/report/thread", feed_errors.MdError(ThreadReport))
	admin.POST("/report/user", feed_errors.MdError(UserReport))
	admin.POST("/trait", feed_errors.MdError(AddCallBlockTrait))

}

type FeedTypeConfForm struct {
	FeedType string `form:"feedType" binding:"required"`
	Conf     string `form:"conf" binding:"required"`
}

// 调用块配置改变
func FeedTypeConfChange(ctx *gin.Context) error {
	var feedTypeConfForm FeedTypeConfForm
	if err := ctx.ShouldBind(&feedTypeConfForm); err != nil {
		return errors.New("params_error")
	}
	if err := api_func.FeedTypeConfChangeService(feedTypeConfForm.FeedType, feedTypeConfForm.Conf); err != nil {
		return errors.New("conf_error")
	}
	ctx.JSON(helper.Success())
	return nil
}

type TopicDataSourceForm struct {
	TopicId  string `form:"topicId" binding:"required"`
	TopicIds string `form:"topicIds" binding:"required"`
}

// topic 数据源改变
func TopicDataSourceChange(ctx *gin.Context) error {
	var topicDataSourceForm TopicDataSourceForm
	if err := ctx.ShouldBind(&topicDataSourceForm); err != nil {
		return errors.New("params_error")
	}
	topicIds := strings.Split(topicDataSourceForm.TopicIds, ",")
	api_func.TopicDataSourceChangeService(topicDataSourceForm.TopicId, topicIds)
	ctx.JSON(helper.Success())
	return nil
}

type TopicForm struct {
	TopicId string `form:"topicId" binding:"required"`
}

// 增加topic
func AddTopic(ctx *gin.Context) error {
	var topicForm TopicForm
	if err := ctx.ShouldBind(&topicForm); err != nil {
		return errors.New("params_error")
	}
	if err := api_func.AddTopicService(topicForm.TopicId); err != nil {
		return err
	}
	ctx.JSON(helper.Success())
	return nil
}

type AgentForm struct {
	TopicId  int    `form:"topicId" binding:"required"`
	FeedType string `form:"feedType" binding:"required"`
	TopicIds string `form:"topicIds" binding:"required"`
}

// 添加agent
func AddAgent(ctx *gin.Context) error {
	var agentForm AgentForm
	if err := ctx.ShouldBind(&agentForm); err != nil {
		return errors.New("params_error")
	}
	topicIds := strings.Split(agentForm.TopicIds, ",")
	if err := api_func.AddAgentService(agentForm.TopicId, agentForm.FeedType, topicIds); err != nil {
		ctx.JSON(helper.Success())
		return nil
	} else {
		return err
	}
}

// 删除topic
func DelTopic(ctx *gin.Context) error {
	var topicForm TopicForm
	if err := ctx.ShouldBind(&topicForm); err != nil {
		return errors.New("params_error")
	}
	api_func.DelTopicService(topicForm.TopicId)
	ctx.JSON(helper.Success())
	return nil
}

// 修改帖子举报规则
func ChangeThreadReportConf(ctx *gin.Context) error {
	var reportThreadConf contract.ReportThreadConf
	if err := ctx.ShouldBind(&reportThreadConf); err != nil {
		return errors.New("params_error")
	}
	api_func.ChangeThreadReportConfService(reportThreadConf)
	ctx.JSON(helper.Success())
	return nil
}

type ThreadReportForm struct {
	ThreadIds string `form:"threadIds" binding:"required"`
}

// 帖子举报
func ThreadReport(ctx *gin.Context) error {
	var threadReportForm ThreadReportForm
	if err := ctx.ShouldBind(&threadReportForm); err != nil {
		return errors.New("params_error")
	}
	tids := strings.Split(threadReportForm.ThreadIds, ",")
	api_func.ThreadReportService(helper.ArrayStrToInt(tids))
	ctx.JSON(helper.Success())
	return nil
}

// 修改用户举报规则
func ChangeUserReportConf(ctx *gin.Context) error {
	var reportUserConf contract.ReportUserConf
	if err := ctx.ShouldBind(&reportUserConf); err != nil {
		return errors.New("params_error")
	}
	api_func.ChangeUserReportConfService(reportUserConf)
	ctx.JSON(helper.Success())
	return nil
}

type UserReportForm struct {
	UserIds string `form:"userIds" binding:"required"`
}

// 用户举报
func UserReport(ctx *gin.Context) error {
	var userReportForm UserReportForm
	if err := ctx.ShouldBind(&userReportForm); err != nil {
		return errors.New("params_error")
	}
	uids := strings.Split(userReportForm.UserIds, ",")
	api_func.UserReportService(helper.ArrayStrToInt(uids))
	ctx.JSON(helper.Success())
	return nil
}

type DelTopicFrom struct {
	TopicId  string `form:"topicId" binding:"required"`
	FeedType string `form:"feedType" binding:"required"`
	Ids      string `form:"ids" binding:"required"`
}

func DelTopicData(ctx *gin.Context) error {
	var delTopicFrom DelTopicFrom
	if err := ctx.ShouldBind(&delTopicFrom); err != nil {
		return errors.New("params_error")
	}
	agentName := fmt.Sprintf("%s%s%s", delTopicFrom.TopicId, service.Separator, delTopicFrom.FeedType)
	ids := strings.Split(delTopicFrom.Ids, ",")
	if err := api_func.DelTopicDataService(agentName, helper.ArrayStrToInt(ids)); err != nil {
		return err
	}
	ctx.JSON(helper.Success())
	return nil
}

type TraitFrom struct {
	Id       string                 `form:"id"`
	TopicId  string                 `form:"topicId" binding:"required"`
	FeedType string                 `form:"feedType" binding:"required"`
	Exp      int                    `form:"exp" binding:"required"`
	Trait    service.CallBlockTrait `form:"trait"`
}

func AddCallBlockTrait(ctx *gin.Context) error {
	var traitFrom TraitFrom
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
