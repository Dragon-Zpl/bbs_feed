package feed_errors

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type errorHandlefunc func(*gin.Context) error

type errorMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var errorMsgs = map[string]errorMsg{
	"params_error": errorMsg{
		Code:    1002,
		Message: "参数解析有误",
	},

	"conf_error": errorMsg{
		Code:    1003,
		Message: "配置传入有误",
	},

	"topic_not_deploy": errorMsg{
		Code:    1004,
		Message: "该topic未配置",
	},
	"redis_key_notexist": errorMsg{
		Code:    1005,
		Message: "redis key 不存在",
	},
}

func GetError(err error) errorMsg {
	if _, ok := errorMsgs[err.Error()]; !ok {
		return errorMsg{
			Code:    1001,
			Message: err.Error(),
		}
	}
	return errorMsgs[err.Error()]
}

func MdError(handlefunc errorHandlefunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err1 := handlefunc(ctx)
		if err1 != nil {
			err2 := GetError(err1)
			ctx.AbortWithStatusJSON(http.StatusOK, err2)
		}
	}
}
