package contribution_list

const (
	Publish_Thread         = "publish_thread"   //发帖
	Thread_Replied         = "thread_replied"   //帖子被回复
	Thread_Collected       = "thread_collected" //帖子被收藏
	Thread_Supported       = "thread_supported" //帖子被加分
	Publish_Thread_Score   = 3                  //发帖得分权重
	Thread_Replied_Score   = 1                  //帖子被回复得分权重
	Thread_Collected_Score = 30                 //帖子被收藏得分权重
	Thread_Supported_Score = 10                 //帖子被加分得分权重
)

func Get() {

}
