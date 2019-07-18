package forms

import "bbs_feed/service"

type TopicForm struct {
	TopicId           string `form:"topicId" binding:"required"`
	HotThread         string `form:"hotThread" binding:"required"`
	Essence           string `form:"essence" binding:"required"`
	TodayIntroduction string `form:"todayIntroduction" binding:"required"`
	WeekPopularity    string `form:"weekPopularity" binding:"required"`
	WeekContribution  string `form:"weekContribution" binding:"required"`
	TopicIds          string `form:"topicIds" binding:"required"`
}

type AgentForm struct {
	TopicId  int    `form:"topicId" binding:"required"`
	FeedType string `form:"feedType" binding:"required"`
	TopicIds string `form:"topicIds" binding:"required"`
}

type FeedTypeConfForm struct {
	FeedType string `form:"feedType" binding:"required"`
	Conf     string `form:"conf" binding:"required"`
}

type TopicDataSourceForm struct {
	TopicId  string `form:"topicId" binding:"required"`
	TopicIds string `form:"topicIds" binding:"required"`
}

type DelTopicFrom struct {
	TopicId  string `form:"topicId" binding:"required"`
	FeedType string `form:"feedType" binding:"required"`
	Ids      string `form:"ids" binding:"required"` //删除指定板块下的tid/uid
}

type ThreadReportForm struct {
	ThreadIds string `form:"threadIds" binding:"required"`
}

type UserReportForm struct {
	UserIds string `form:"userIds" binding:"required"`
}

type TraitFrom struct {
	Id       string                 `form:"id"`
	TopicId  string                 `form:"topicId" binding:"required"`
	FeedType string                 `form:"feedType" binding:"required"`
	Exp      int                    `form:"exp" binding:"required"`
	Trait    service.CallBlockTrait `form:"trait"`
}


type CallBackArgs struct {
	TopicId string `form:"topicId" binding:"required"`
	Block string `form:"block" binding:"required"`
}

//type ThreadStruct struct {
//	Tid int `json:"tid"`
//	Fid int `json:"fid"`
//	Author string `json:"author"`
//	AuthorId int `json:"authorId"`
//	Subject string `json:"subject"`
//	Dateline int `json:"dateline"`
//}
//
//type TraitStruct struct {
//	IsSetTop bool `json:"isSetTop"`
//	Subscript string `json:"subscript"`
//}
//
//type RedisEssence struct {
//	Thread ThreadStruct `json:"thread"`
//	Trait TraitStruct `json:"trait"`
//}