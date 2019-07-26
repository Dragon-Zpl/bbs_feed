package contribution_list

import (
	"bbs_feed/lib/helper"
	"bbs_feed/search"
	"github.com/astaxie/beego/logs"
)

const (
	Publish_Thread_Score   = 3                  //发帖得分权重
	Thread_Replied_Score   = 1                  //帖子被回复得分权重
	Thread_Collected_Score = 30                 //帖子被收藏得分权重
	Thread_Supported_Score = 10                 //帖子被加分得分权重
)



func GetActionScore(useraAction search.UserAction) int {
	return useraAction.PublishThread * Publish_Thread_Score + useraAction.ThreadReplied * Thread_Replied_Score + useraAction.ThreadCollected * Thread_Collected_Score + useraAction.ThreadSupported * Thread_Supported_Score
}

func GetContributeScore() (map[string][]*search.User, error) {
	index := helper.GetWeekStart().Format("2006-01-02")
	esDatas, err := search.Search(index)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	for _, v := range esDatas{
		for _, data := range v{
			data.Score = GetActionScore(data.Action)
		}
	}
	return esDatas, nil
}
