package popularity_list

import (
	"bbs_feed/lib/helper"
	"bbs_feed/search"
	"github.com/astaxie/beego/logs"
)

func GetPopularityScore() (map[string][]*search.User, error) {
	index := helper.GetWeekStart().Format("2006-01-02")
	esDatas, err := search.Search(index)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	return esDatas, nil
}
