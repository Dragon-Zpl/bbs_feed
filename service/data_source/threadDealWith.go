package data_source

import (
	"bbs_feed/boot"
	"bbs_feed/model/forum_thread"
	"bbs_feed/service"
	"bbs_feed/service/redis_ops"
	"encoding/json"
	"sort"
	"strconv"
)

type Threads []*forum_thread.Model

type RedisThread struct {
	Thread forum_thread.Model     `json:"thread"`
	Trait  service.CallBlockTrait `json:"trait"`
}

func (this Threads) Len() int {
	return len(this)
}

func (this Threads) Less(i, j int) bool {
	if this[i].Replies < this[j].Replies {
		return true
	} else if this[i].Replies == this[j].Replies {
		if this[i].FavTimes < this[j].FavTimes {
			return true
		} else if this[i].FavTimes == this[j].FavTimes {
			if this[i].Dateline < this[j].Dateline {
				return true
			}
		}
	}
	return false
}

func (this Threads) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

// 获取排序后的thread
func GetHotSortThread(fids []int, day, views, replys int) []RedisThread {
	var (
		hotThreads   Threads
		redisThreads []RedisThread
	)
	hotThreads = forum_thread.GetHotThreads(fids, day, views, replys)
	redisThreads = make([]RedisThread, 0, len(hotThreads))

	sort.Sort(hotThreads)
	for _, thread := range hotThreads {
		redisThreads = append(redisThreads, RedisThread{
			Thread: *thread,
		})
	}
	return redisThreads
}

func GetEssenceSortThread(fids []int, day int) []RedisThread {
	var (
		essEnceThreads Threads
		redisThreads   []RedisThread
	)
	essEnceThreads = forum_thread.GetEssenceThreads(fids, day)
	redisThreads = make([]RedisThread, 0, len(essEnceThreads))

	sort.Sort(essEnceThreads)
	for _, thread := range essEnceThreads {
		redisThreads = append(redisThreads, RedisThread{
			Thread: *thread,
		})
	}
	return redisThreads
}

func DelRedisThreadInfo(tids []int, key, traitKey string) {
	essenceThreads := GetThreadByTids(tids)
	for _, thread := range essenceThreads {
		if trait, err := boot.InstanceRedisCli(boot.CACHE).HGet(traitKey, strconv.Itoa(thread.Thread.Tid)).Result(); err == nil {
			var callBlockTrait service.CallBlockTrait
			if err = json.Unmarshal([]byte(trait), &callBlockTrait); err == nil {
				thread.Trait = callBlockTrait
			}
		}
		if threadBytes, err := json.Marshal(thread); err == nil {
			redis_ops.DelZAdd(key, string(threadBytes))
		}
	}
}

func GetThreadByTids(tids []int) []RedisThread {
	var redisThreads []RedisThread
	threads := forum_thread.GetByTids(tids)
	for _, thread := range threads {
		redisThreads = append(redisThreads, RedisThread{
			Thread: *thread,
		})
	}
	return redisThreads
}
