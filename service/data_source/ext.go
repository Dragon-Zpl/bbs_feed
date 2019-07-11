package data_source

import (
	"bbs_feed/model/topic"
	"strconv"
)

// 获取全局的fids
func GetAllTopics() []string {
	allTopics := topic.GetAll()
	fids := make([]string, 0, len(allTopics))
	for _, t := range allTopics {
		fids = append(fids, strconv.Itoa(t.Id))
	}
	return fids
}


// 检查是否需要全局的fids  ["0"] 代表全局
func CheckIsGetAllTopics(inTopics []string) bool {
	if len(inTopics) == 1 && inTopics[0] == "0"{
		return true
	}
	return false
}

// 获取正确的topics
func CheckAngGetCurTopics(inTopics []string) []string {
	if CheckIsGetAllTopics(inTopics) {
		return GetAllTopics()
	} else {
		return inTopics
	}
}