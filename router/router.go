package router

import (
	"bbs_feed/v1/admin"
	"bbs_feed/v1/view"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	var app = gin.New()
	app.Use(gin.Recovery())
	admin.Mapping("/admin", app)
	view.Mapping("/view", app)
	return app
}
