package schedules

import (
	"bbs_feed/conf"
	"bbs_feed/search"
	"github.com/astaxie/beego/logs"
	"github.com/robfig/cron"
	"time"
)

const (
	//定时删除上上周 周一es中人气榜、贡献榜数据的索引
	del_spec = "30 0 1 * * 1"
)

func Crontab() {
	delIndex()
}

func delIndex() {
	c := cron.New()
	c.AddFunc(del_spec, func() {
		twoWeeksAgo := time.Now().AddDate(0, 0, -14).Format("2006-01-02")
		index := conf.EsConf.Index + "_" + twoWeeksAgo
		if err := search.DeleteIndex(index); err != nil {
			logs.Error(err)
		}
	})
	c.Start()
}
