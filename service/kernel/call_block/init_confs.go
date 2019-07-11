package call_block

import (
	"bbs_feed/model/feed_conf"
	"encoding/json"
)

/*
读取各个模块的配置
*/

var (
	hot            HotRules
	essence        EssenceRules
	contribution   ContributionRules
	weekPopularity WeekPopularityRule
	newHot         NewHotRules
	todayIntro     IntroRules
)

func InitConfs() {
	confs := feed_conf.GetAll()
	for _, conf := range confs {
		switch conf.Name {
		case "hot":
			json.Unmarshal([]byte(conf.Conf), &hot)
		case "essence":
			json.Unmarshal([]byte(conf.Conf), &essence)
		case "contribution":
			json.Unmarshal([]byte(conf.Conf), &contribution)
		case "weekPopularity":
			json.Unmarshal([]byte(conf.Conf), &weekPopularity)
		case "hot_new":
			json.Unmarshal([]byte(conf.Conf), &newHot)
		case "todayIntro":
			json.Unmarshal([]byte(conf.Conf), &todayIntro)
		}
	}
}
