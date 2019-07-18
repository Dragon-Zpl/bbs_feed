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
<<<<<<< HEAD
	//creater.InitService()
=======
	creater.InitService()
>>>>>>> a09d64169b591812f8d56a69765ffaa87abb171a
	r := router.Router()
	r.Run(":8887")
}
