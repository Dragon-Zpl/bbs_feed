package view

import (
	"bbs_feed/lib/feed_errors"
	"bbs_feed/lib/helper"
	"bbs_feed/model/feed_permission"
	"errors"
	"github.com/gin-gonic/gin"
)

func Mapping(prefix string, app *gin.Engine)  {
	admin := app.Group(prefix)
	admin.GET("/call_back", feed_errors.MdError(GetCallBack))
}

type CallBackArgs struct {
	TopicId string `form:"topicId" binding:"required"`
	Block string `form:"block" binding:"required"`
}

func GetCallBack(ctx *gin.Context) error{
	var callBackArgs CallBackArgs
	if err := ctx.ShouldBindQuery(&callBackArgs) ; err != nil{
		return errors.New("params_error")
	}
	if err := feed_permission.GetBlock(callBackArgs.TopicId, callBackArgs.Block) ; err==nil{
		ctx.JSON(helper.Success())
	} else {
		return errors.New("topic_not_deploy")
	}
	//ctx.JSON(helper.SuccessWithDate(callBackArgs))
	return nil
}