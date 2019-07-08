package data_source

import (
	"bbs_feed/model/forum_thread"
	"bbs_feed/service"
	"sort"
)

type HotThreads []*forum_thread.Model

type HotRedisThread struct {
	Thread forum_thread.Model `json:"thread"`
	Trait service.CallBlockTrait `json:"trait"`
}


func (this HotThreads) Len() int {
	return len(this)
}

func (this HotThreads) Less(i, j int) bool {
	if this[i].Replies < this[j].Replies {
		return true
	} else if this[i].Replies == this[j].Replies {
		if this[i].FavTimes < this[j].FavTimes {
			return true
		} else if this[i].FavTimes == this[j].FavTimes {
			if this[i].Dateline	 < this[j].Dateline {
				return true
			}
		}
	}
	return false
}

func (this HotThreads) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

// 获取排序后的thread
func GetHotSortThread(fids []int, day, views, replys int) []HotRedisThread {
	var (
		hotThreads HotThreads
		redisThreads []HotRedisThread
	)
	hotThreads = forum_thread.GetByCondition(fids, day, views, replys)
	redisThreads = make([]HotRedisThread, 0, len(hotThreads))

	sort.Sort(hotThreads)
	for _, thread := range hotThreads {
		redisThreads = append(redisThreads, HotRedisThread{
			Thread:          *thread,
		})
	}
	return redisThreads
}

func GetHotThreadByTids(tids []int)[]HotRedisThread  {
	var redisThreads []HotRedisThread
	threads := forum_thread.GetByTids(tids)
	for _, thread := range threads {
		redisThreads = append(redisThreads, HotRedisThread{
			Thread:          *thread,
		})
	}
	return redisThreads
}