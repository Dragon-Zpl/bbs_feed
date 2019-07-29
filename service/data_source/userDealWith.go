package data_source

import (
	"bbs_feed/model/common_member"
	"bbs_feed/model/contribution_list"
	"bbs_feed/model/popularity_list"
	"bbs_feed/search"
	"bbs_feed/service"
	"github.com/astaxie/beego/logs"
	"sort"
	"strconv"
)

type User struct {
	Uid   int    `json:"uid"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type RedisUser struct {
	User  User                   `json:"user"`
	Trait service.CallBlockTrait `json:"trait"`
}

func GetUserActionAdd(ids []int, populData map[string][]*search.User) map[string]*search.UserAction {
	actionSlice := make(map[string]*search.UserAction)
	for _, id := range ids {
		if _, ok := populData[strconv.Itoa(id)]; !ok {
			continue
		}
		for _, data := range populData[strconv.Itoa(id)] {
			if _, ok := actionSlice[data.Uid]; !ok {
				actionSlice[data.Uid] = &data.Action
				continue
			}
			actionSlice[data.Uid].ThreadSupported += data.Action.ThreadSupported
			actionSlice[data.Uid].ThreadCollected += data.Action.ThreadCollected
			actionSlice[data.Uid].ThreadReplied += data.Action.ThreadReplied
			actionSlice[data.Uid].PublishThread += data.Action.PublishThread
			actionSlice[data.Uid].PostSupported += data.Action.PostSupported
			actionSlice[data.Uid].PublishPost += data.Action.PublishPost
		}
	}
	return actionSlice
}

func GetPopulData(ids []int, limit, PostSupportedScore, PublishPostScore, PublishThreadScore, ThreadSupportedScore int) []*RedisUser {
	populData, err := popularity_list.GetPopularityScore()
	if err != nil {
		logs.Error(err)
	}
	resSlice := make(PopulAndContriSort, 0)
	var name string
	actionSlice := GetUserActionAdd(ids, populData)
	for id, userAction := range actionSlice {
		redisUser := new(RedisUser)
		user := new(User)
		uid, _ := strconv.Atoi(id)
		user.Uid = uid
		user.Score = userAction.PublishPost*PublishPostScore + userAction.PostSupported*PostSupportedScore + userAction.PublishThread*PublishThreadScore + userAction.ThreadSupported*ThreadSupportedScore
		userName, err := common_member.GetUserName(id)
		if err != nil {
			name = ""
		} else {
			name = userName.Username
		}
		user.Name = name
		redisUser.User = *user
		resSlice = append(resSlice, redisUser)
	}

	sort.Sort(resSlice)

	if len(resSlice) < limit {
		return resSlice
	}

	return resSlice[:limit]
}

func GetContributionData(ids []int, limit, Publish_Thread_Score, Thread_Replied_Score, Thread_Collected_Score, Thread_Supported_Score int) []*RedisUser {
	contributionData, err := contribution_list.GetContributeScore()
	if err != nil {
		logs.Error(err)
	}
	resSlice := make(PopulAndContriSort, 0)
	var name string
	actionSlice := GetUserActionAdd(ids, contributionData)
	for id, userAction := range actionSlice {
		redisUser := new(RedisUser)
		user := new(User)
		uid, _ := strconv.Atoi(id)
		user.Uid = uid
		user.Score = userAction.PublishPost*Publish_Thread_Score + userAction.PostSupported*Thread_Replied_Score + userAction.PublishThread*Thread_Collected_Score + userAction.ThreadSupported*Thread_Supported_Score
		userName, err := common_member.GetUserName(id)
		if err != nil {
			name = ""
		} else {
			name = userName.Username
		}
		user.Name = name
		redisUser.User = *user
		resSlice = append(resSlice, redisUser)
	}

	sort.Sort(resSlice)
	if len(resSlice) < limit {
		return resSlice
	}

	return resSlice[:limit]
}
