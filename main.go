package main

import (
	"bbs_feed/boot"
	"bbs_feed/conf"
	"bbs_feed/router"
	"bbs_feed/schedules"
	"bbs_feed/service/kernel/creater"
)

func main() {
	conf.InitConf()
	boot.ConnectMySQL()
	boot.ConnectRedis()
	boot.InitSearchClient()
	creater.InitService()
	go schedules.Crontab()
	r := router.Router()
	r.Run(":8887")
}
