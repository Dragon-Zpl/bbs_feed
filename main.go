package main

import (
	"bbs_feed/boot"
	"bbs_feed/conf"
	"bbs_feed/router"
	"bbs_feed/service/kernel/creater"
)

func main() {
	conf.InitConf()
	boot.ConnectMySQL()
	creater.InitService()
	r := router.Router()
	r.Run(":8887")
}
