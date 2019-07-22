#!/usr/bin/env bash

#IP
feed="root@192.168.188.84"
basepath=$(cd `dirname $0`; pwd)


if [ $1 == "dev" ]
then
    #编译
    echo "编译中..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bbs_feed main.go

    #删除服务器上的bbs_feed
    echo "rm -rf bbs_feed..."
    ssh -p 58422 $feed << feed
        rm -rf /www/bbs_feed/bbs_feed
        exit
feed

    #部署
    echo "上传bbs_feed..."
    scp -P 58422 $basepath/bbs_feed $feed:/www/bbs_feed/
    scp -P 58422 $basepath/conf/conf.ini $feed:/www/bbs_feed/conf/

    #重启服务
    echo "pm2 reload bbs_feed..."
    ssh -p 58422 $feed << feed
        pm2 reload /www/bbs_feed/bbs_feed
        exit
feed

else
    echo "缺失部署环境参数！！！"
fi