package service_confs

import (
	"bbs_feed/model/feed_conf"
	"bbs_feed/service/kernel"
	"encoding/json"
)

var (
	Hot kernel.HotRules
	Essence kernel.EssenceRules
	Contribution kernel.ContributionRules
	WeekPopularity kernel.WeekPopularityRule
)



func InitConfs() {
	confs := feed_conf.GetAll()
	for _, conf := range confs {
		switch conf.Name {
		case "hot":
			json.Unmarshal([]byte(conf.Conf), &Hot)
		case "essence":
			json.Unmarshal([]byte(conf.Conf), &Essence)
		case "contribution":
			json.Unmarshal([]byte(conf.Conf), &Contribution)
		case "weekPopularity":
			json.Unmarshal([]byte(conf.Conf), &WeekPopularity)

		}
	}
}