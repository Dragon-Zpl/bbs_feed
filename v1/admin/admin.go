package admin

import (
	"bbs_feed/lib/helper"
	"bbs_feed/service"
	"bbs_feed/service/kernel/contract"
	"bbs_feed/service/kernel/creater"
	"bbs_feed/service/redis_ops"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func Mapping(prefix string, app *gin.Engine) {
	admin := app.Group(prefix)
	admin.POST("/topic", AddTopic)
	admin.POST("/agnet", AddAgent)
	admin.DELETE("/topic/:id", DelTopic)
	admin.POST("/topic/conf", FeedTypeConfChange)
	admin.POST("/topic/fids", TopicDataSourceChange)
	admin.POST("/report/thread/conf", ChangeThreadReportConf)
	admin.POST("/report/user/conf", ChangeUserReportConf)
	admin.DELETE("/call_back", DelTopicData)
	admin.POST("/report/thread", ThreadReport)
	admin.POST("/report/user", UserReport)
	admin.POST("/trait", AddCallBlockTrait)

}

type FeedTypeConfForm struct {
	FeedType string `form:"feedType" binding:"required"`
	Conf     string `form:"conf" binding:"required"`
}

// 调用块配置改变
func FeedTypeConfChange(ctx *gin.Context) {
	var feedTypeConfForm FeedTypeConfForm
	if err := ctx.ShouldBind(&feedTypeConfForm); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":    1053,
			"message": "参数传入有误",
		})
		return
	}
	if err := contract.InstanceFeedService().ChangeConf(feedTypeConfForm.FeedType, feedTypeConfForm.Conf); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":    1053,
			"message": "配置传入有误",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

type TopicDataSourceForm struct {
	TopicId  string `form:"topicId" binding:"required"`
	TopicIds string `form:"topicIds" binding:"required"`
}

// topic 数据源改变
func TopicDataSourceChange(ctx *gin.Context) {
	var topicDataSourceForm TopicDataSourceForm
	if err := ctx.ShouldBind(&topicDataSourceForm); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":    1053,
			"message": "参数传入有误",
		})
		return
	}
	topicIds := strings.Split(topicDataSourceForm.TopicIds, ",")
	contract.InstanceFeedService().ChangeFids(topicDataSourceForm.TopicId, topicIds)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

type TopicForm struct {
	TopicId string `form:"topicId" binding:"required"`
}

// 增加topic
func AddTopic(ctx *gin.Context) {
	var topicForm TopicForm
	if err := ctx.ShouldBind(&topicForm); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":    1053,
			"message": "参数传入有误",
		})
		return
	}
	if agents, err := creater.GenAgents(topicForm.TopicId); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":    1053,
			"message": fmt.Sprintf("topic %s 未配置", topicForm.TopicId),
		})
		return
	} else {
		contract.InstanceFeedService().RegisterService(agents...)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

type AgentForm struct {
	TopicId int `form:"topicId" binding:"required"`
	FeedType string `form:"feedType" binding:"required"`
	TopicIds string `form:"topicIds" binding:"required"`
}

// 添加agent
func AddAgent(ctx *gin.Context)  {
	var agentForm AgentForm
	if err := ctx.ShouldBind(&agentForm); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":    1053,
			"message": "参数传入有误",
		})
		return
	}
	topicIds := strings.Split(agentForm.TopicIds, ",")
	if agent := creater.GenAgent(agentForm.TopicId, agentForm.FeedType, topicIds); agent == nil {
		contract.InstanceFeedService().RegisterService(agent)
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
		})
		return
	} else {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":    1053,
			"message": "参数传入有误",
		})
		return
	}
}

// 删除topic
func DelTopic(ctx *gin.Context) {
	var topicForm TopicForm
	if err := ctx.ShouldBind(&topicForm); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":    1053,
			"message": "参数传入有误",
		})
		return
	}
	contract.InstanceFeedService().RemovePusher(topicForm.TopicId)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// 修改帖子举报规则
func ChangeThreadReportConf(ctx *gin.Context) {
	var reportThreadConf contract.ReportThreadConf
	if err := ctx.ShouldBind(&reportThreadConf); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":    1053,
			"message": "参数传入有误",
		})
		return
	}
	contract.ThreadReportCheck.ChangeConf(reportThreadConf)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

type ThreadReportForm struct {
	ThreadIds string `form:"threadIds" binding:"required"`
}

// 帖子举报
func ThreadReport(ctx *gin.Context) {
	var threadReportForm ThreadReportForm
	if err := ctx.ShouldBind(&threadReportForm); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":    1053,
			"message": "参数传入有误",
		})
		return
	}
	tids := strings.Split(threadReportForm.ThreadIds, ",")
	tidsInt := helper.ArrayStrToInt(tids)
	contract.ThreadReportCheck.AcceptReportTids(tidsInt)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// 修改用户举报规则
func ChangeUserReportConf(ctx *gin.Context) {
	var reportUserConf contract.ReportUserConf
	if err := ctx.ShouldBind(&reportUserConf); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":    1053,
			"message": "参数传入有误",
		})
		return
	}
	contract.UserReportCheck.ChangeConf(reportUserConf)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

type UserReportForm struct {
	UserIds string `form:"userIds" binding:"required"`
}

// 用户举报
func UserReport(ctx *gin.Context) {
	var userReportForm UserReportForm
	if err := ctx.ShouldBind(&userReportForm); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":    1053,
			"message": "参数传入有误",
		})
		return
	}
	uids := strings.Split(userReportForm.UserIds, ",")
	uidInts := helper.ArrayStrToInt(uids)
	contract.UserReportCheck.AcceptReportUids(uidInts)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

type DelTopicFrom struct {
	TopicId  string `form:"topicId" binding:"required"`
	FeedType string `form:"feedType" binding:"required"`
	Ids      string `form:"ids" binding:"required"`
}

func DelTopicData(ctx *gin.Context) {
	var delTopicFrom DelTopicFrom
	if err := ctx.ShouldBind(&delTopicFrom); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":    1053,
			"message": "参数传入有误",
		})
		return
	}
	agentName := fmt.Sprintf("%s%s%s", delTopicFrom.TopicId, service.Separator, delTopicFrom.FeedType)
	ids := strings.Split(delTopicFrom.Ids, ",")
	idsInt := helper.ArrayStrToInt(ids)
	if err := contract.InstanceFeedService().Remove(agentName, idsInt); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":    1053,
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

type TraitFrom struct {
	Id       string                 `form:"id"`
	TopicId  string                 `form:"topicId" binding:"required"`
	FeedType string                 `form:"feedType" binding:"required"`
	Exp      int `form:"exp" binding:"required"`
	Trait    service.CallBlockTrait `form:"trait"`
}

func AddCallBlockTrait(ctx *gin.Context) {
	var traitFrom TraitFrom
	if err := ctx.ShouldBind(&traitFrom); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":    1053,
			"message": "参数传入有误",
		})
		return
	}
	redisKey := fmt.Sprintf("call_block_%s_trait", traitFrom.FeedType)
	if !redis_ops.KeyExist(redisKey) {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":    1053,
			"message": "参数传入有误",
		})
		return
	}
	if traitBytes, err := json.Marshal(traitFrom.Trait); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":    1053,
			"message": "参数传入有误",
		})
		return
		redis_ops.HSet(redisKey, traitFrom.Id, string(traitBytes), time.Duration(traitFrom.Exp) * time.Hour)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
