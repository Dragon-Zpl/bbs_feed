#!/usr/bin/env bash

feed="root@192.168.188.84"
basepath=$(
	cd $(dirname $0)
	pwd
)

deploy() {
	echo "rm -rf bbs_feed..."
	ssh -p 58422 $feed <<ssh
        rm -rf /www/bbs_feed/bbs_feed
        exit
ssh
	echo "uploading...[192.168.188.84]"
	scp -P 58422 $basepath/bbs_feed $feed:/www/bbs_feed/
	scp -P 58422 $basepath/conf/conf.ini $feed:/www/bbs_feed/conf/
	ssh -p 58422 $feed <<ssh
        pm2 reload /www/bbs_feed/bbs_feed
        exit
ssh
}

echo "compiling..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bbs_feed main.go

if [ $1 == "dev" ]; then
	deploy
else
	echo "missing deployment parameters..."
fi
