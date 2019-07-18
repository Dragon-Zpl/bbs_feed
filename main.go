package main

import (
	"bbs_feed/boot"
	"bbs_feed/conf"
	"bbs_feed/router"
)

func main() {
	conf.InitConf()
	boot.ConnectMySQL()
	boot.ConnectRedis()
	//creater.InitService()
	r := router.Router()
	r.Run(":8887")
}
