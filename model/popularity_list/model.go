package popularity_list

import (
	"bbs_feed/lib/helper"
	"bbs_feed/search"
	"fmt"
)

const (
	Thread_Supported_Score = 10                 //帖子被加分权重
	Publish_Thread_Score   = 3                  //发帖权重
	Publish_Post_Score     = 1                  //评论权重
	Post_Supported_Score   = 5                  //评论被加分权重
	Thread_Supported       = "thread_supported" //帖子被加分
	Publish_Thread         = "publish_thread"   //发帖
	Publish_Post           = "publish_post"     //评论
	Post_Supported         = "post_supported"   //评论被加分
)

func GetActionScore(useraAction search.UserAction) int {
	return useraAction.ThreadSupported * 10 + useraAction.PublishThread * 3 + useraAction.PublishPost + useraAction.PostSupported * 5
}

func GetUserAction() {
	index := helper.GetWeekStart().Format("2006-01-02")
	esDatas, err := search.Search(index)
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range esDatas{
		for _, data := range v{
			data.Score = GetActionScore(data.Action)
		}
	}
	fmt.Println(esDatas)
}
