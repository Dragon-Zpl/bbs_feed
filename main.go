package main

import (
	"bbs_feed/boot"
	"bbs_feed/conf"
	"bbs_feed/service/kernel"
	"bbs_feed/service/kernel/call_block"
)

func main() {
	conf.InitConf()
	boot.ConnectMySQL()
	kernel.NewSerivce()
	call_block.Remove([]int{12372363})
	//time.Sleep(10 * time.Minute)
}

