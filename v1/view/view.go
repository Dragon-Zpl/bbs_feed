package view

import (
	"bbs_feed/lib/feed_errors"
	"bbs_feed/lib/helper"
	"bbs_feed/model/feed_conf"
	"bbs_feed/service/api_func"
	"bbs_feed/v1/forms"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

var Blocks = map[string]int{"hot_thread" : 1, "essence" : 1, "today_introduction" : 1, "week_popularity" : 1, "week_contribution" : 1}

func Mapping(prefix string, app *gin.Engine)  {
	admin := app.Group(prefix)
	admin.GET("/call_back", feed_errors.MdError(GetCallBack))
	admin.GET("/feed_conf", feed_errors.MdError(GetFeedConf))
}

func GetFeedConf(ctx *gin.Context) error{
	var datas []*feed_conf.Model
	datas = feed_conf.GetAll()
	if datas != nil{
		for _, data := range datas{
			data.Conf = strings.Replace(data.Conf,"\"","'", -1)
		}
		ctx.JSON(helper.SuccessWithDate(datas))
	}
	return nil
}

func GetCallBack(ctx *gin.Context) error{
	var callBackArgs forms.CallBackArgs
	if err := ctx.ShouldBindQuery(&callBackArgs) ; err != nil{
		return errors.New("params_error")
	}
	if _,ok := Blocks[callBackArgs.Block] ; !ok{
		return errors.New("params_error")
	}
	if res_data, err := api_func.GetRedisBlockDataService(callBackArgs.TopicId, callBackArgs.Block) ; err == nil{
		ctx.JSON(helper.SuccessWithDate(res_data))
	} else{
		return err
	}
	return nil
}