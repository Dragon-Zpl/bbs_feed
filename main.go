package main

import (
	"bbs_feed/boot"
	"bbs_feed/conf"
	"bbs_feed/router"
	"bbs_feed/service/kernel/contract"
)

func main() {
	conf.InitConf()
	boot.ConnectMySQL()
	contract.InitFeedService()
	contract.NewThreadReportCheckEr()
	contract.NewUserReportCheck()
	r := router.Router()
	r.Run(":8888")
}
