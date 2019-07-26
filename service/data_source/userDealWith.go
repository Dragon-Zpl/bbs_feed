package data_source

import (
	"bbs_feed/model/common_member"
	"bbs_feed/model/contribution_list"
	"bbs_feed/model/popularity_list"
	"bbs_feed/model/topic_fid_relation"
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

func GetPopulData(limit int) map[string][]RedisUser {
	populData, err := popularity_list.GetPopularityScore()
	if err != nil {
		logs.Error(err)
	}
	fidTopicId, err := topic_fid_relation.GetTopicIds()
	var name string
	redisData := make(map[string][]RedisUser)
	for k, v := range populData {
		var populSort PopulAndContriSort
		populSort = v
		sort.Sort(populSort)
		topicId := fidTopicId[k]
		limitData := v
		// 需要修改
		if len(v) > limit {
			limitData = v[:limit]
		}
		if _, ok := redisData[topicId]; !ok {
			redisData[topicId] = make([]RedisUser, 0)
		}

		for _, data := range limitData {
			redisUser := new(RedisUser)
			user := new(User)
			uid, _ := strconv.Atoi(data.Uid)
			user.Uid = uid
			user.Score = data.Score
			userName, err := common_member.GetUserName(data.Uid)
			if err != nil {
				name = ""
			} else {
				name = userName.Username
			}
			user.Name = name
			redisUser.User = *user
			redisData[topicId] = append(redisData[topicId], *redisUser)
		}
	}
	return redisData
}

func GetContributionData(limit int) map[string][]RedisUser {
	contributionData, err := contribution_list.GetContributeScore()
	if err != nil {
		logs.Error(err)
	}
	fidTopicId, err := topic_fid_relation.GetTopicIds()
	var name string
	redisData := make(map[string][]RedisUser)
	for k, v := range contributionData {
		var populSort PopulAndContriSort
		populSort = v
		sort.Sort(populSort)
		topicId := fidTopicId[k]
		limitData := v
		// 需要修改
		if len(v) > limit {
			limitData = v[:limit]
		}
		if _, ok := redisData[topicId]; !ok {
			redisData[topicId] = make([]RedisUser, 0)
		}

		for _, data := range limitData {
			redisUser := new(RedisUser)
			user := new(User)
			uid, _ := strconv.Atoi(data.Uid)
			user.Uid = uid
			user.Score = data.Score
			userName, err := common_member.GetUserName(data.Uid)
			if err != nil {
				name = ""
			} else {
				name = userName.Username
			}
			user.Name = name
			redisUser.User = *user
			redisData[topicId] = append(redisData[topicId], *redisUser)
		}
	}
	return redisData
}
