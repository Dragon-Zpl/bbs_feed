package view

import (
	"bbs_feed/lib/feed_errors"
	"bbs_feed/lib/helper"
	"bbs_feed/model/feed_conf"
	"bbs_feed/model/feed_permission"
	"bbs_feed/service/api_func"
	"github.com/gin-gonic/gin"
	"strings"
)

func Mapping(prefix string, app *gin.Engine) {
	admin := app.Group(prefix)
	admin.GET("/feed_conf", feed_errors.MdError(GetFeedConf))
	admin.GET("/feed_premission", feed_errors.MdError(GetPremission))
	admin.GET("/feed_conf_use", feed_errors.MdError(GetFeedConfUse))
	admin.GET("/topic", feed_errors.MdError(GetTopic))
}

func GetFeedConf(ctx *gin.Context) error {
	var datas []*feed_conf.Model
	datas = feed_conf.GetAll()
	if datas != nil {
		for _, data := range datas {
			data.Conf = strings.Replace(data.Conf, "\"", "'", -1)
		}
		ctx.JSON(helper.SuccessWithDate(datas))
	}
	return nil
}

func GetFeedConfUse(ctx *gin.Context) error {
	block_datas := api_func.GetFeedConfUseSerive()
	ctx.JSON(helper.SuccessWithDate(block_datas))
	return nil
}

func GetPremission(ctx *gin.Context) error {
	datas := feed_permission.GetAll()
	ctx.JSON(helper.SuccessWithDate(datas))
	return nil
}

func GetTopic(ctx *gin.Context) error {
	res_datas := api_func.GetTopicSerive()
	ctx.JSON(helper.SuccessWithDataList(res_datas))
	return nil
}
