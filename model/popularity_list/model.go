package popularity_list

import (
	"bbs_feed/lib/helper"
	"bbs_feed/search"
	"github.com/astaxie/beego/logs"
)

const (
	Thread_Supported_Score = 10 //帖子被加分权重
	Publish_Thread_Score   = 3  //发帖权重
	Publish_Post_Score     = 1  //评论权重
	Post_Supported_Score   = 5  //评论被加分权重
)

func GetActionScore(useraAction search.UserAction) int {
	return useraAction.ThreadSupported*Thread_Supported_Score + useraAction.PublishThread*Publish_Thread_Score + useraAction.PublishPost*Publish_Post_Score + useraAction.PostSupported*Post_Supported_Score
}

func GetPopularityScore() (map[string][]*search.User, error) {
	index := helper.GetWeekStart().Format("2006-01-02")
	esDatas, err := search.Search(index)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	for _, v := range esDatas {
		for _, data := range v {
			data.Score = GetActionScore(data.Action)
		}
	}
	return esDatas, nil
}
