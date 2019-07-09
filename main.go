package main

import (
	"bbs_feed/boot"
	"bbs_feed/conf"
	"bbs_feed/router"
	"bbs_feed/service/kernel/contract"
	"bbs_feed/service/kernel/creater"
)

func main() {
	conf.InitConf()
	boot.ConnectMySQL()
	contract.NewFeedService(creater.CreateAgents()...)
	contract.NewThreadReportCheckEr()
	contract.NewUserReportCheck()
	r := router.Router()
	r.Run(":8888")
}
