package admin

import (
	"bbs_feed/service/kernel"
	"github.com/gin-gonic/gin"
)

// 调用块配置改变
func FeedTypeConfChange(ctx *gin.Context) {
	kernel.FeedService.ChangeConf()
}

// topic 数据源改变
func TopicDataSouceChange(ctx *gin.Context) {
	kernel.FeedService.ChangeFids()
}

// 增加topic
func AddTopic(ctx *gin.Context) {
	kernel.FeedService.RegisterService()
}

// 删除topic
func DelTopic(ctx *gin.Context) {
	kernel.FeedService.StopPusher()
}

// 修改帖子举报规则
func ChangeThreadReportConf(ctx *gin.Context) {
	kernel.ThreadReportCheck.ChangeConf()
}

// 修改用户举报规则
func ChangeUserReportConf(ctx *gin.Context) {
	kernel.UserReportCheck.ChangeConf()
}