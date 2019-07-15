package data_source

import (
	"bbs_feed/model/forum_thread"
	"bbs_feed/service"
	"bbs_feed/service/redis_ops"
	"encoding/json"
	"sort"
)

type RedisThread struct {
	Thread forum_thread.Model     `json:"thread"`
	Trait  service.CallBlockTrait `json:"trait"`
}

// 获取hot排序后的thread
func GetHotSortThread(fids []int, day, views, replys int) []RedisThread {
	var (
		hotThreads   HotThread
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

// 获取 精华 排序后的thread
func GetEssenceSortThread(fids []int, day int) []RedisThread {
	var (
		essEnceThreads EssenceThread
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

// 获取今日导读排序后的thread
func GetTodayIntroSortThread(fids []int, day int) []RedisThread {
	var (
		todayIntroThreads TodayIntroThread
		redisThreads      []RedisThread
	)
	todayIntroThreads = forum_thread.GetHotThreads(fids, day, 0, 0)
	redisThreads = make([]RedisThread, 0, len(todayIntroThreads))

	sort.Sort(todayIntroThreads)
	for _, thread := range todayIntroThreads {
		redisThreads = append(redisThreads, RedisThread{
			Thread: *thread,
		})
	}
	return redisThreads
}

// 获取今日导读排序后的thread
func GetNewHotSortThread(fids []int, day int) []RedisThread {
	var (
		newHotThreads NewHotThread
		redisThreads  []RedisThread
	)
	newHotThreads = forum_thread.GetHotThreads(fids, day, 0, 0)
	redisThreads = make([]RedisThread, 0, len(newHotThreads))

	sort.Sort(newHotThreads)
	for _, thread := range newHotThreads {
		redisThreads = append(redisThreads, RedisThread{
			Thread: *thread,
		})
	}
	return redisThreads
}

func DelRedisThreadInfo(tids []int, key, traitKey string) {
	threads := GetThreadByTids(tids)
	for _, thread := range threads {
		//if trait, err := boot.InstanceRedisCli(boot.CACHE).HGet(traitKey, strconv.Itoa(thread.Thread.Tid)).Result(); err == nil {
		//	var callBlockTrait service.CallBlockTrait
		//	if err = json.Unmarshal([]byte(trait), &callBlockTrait); err == nil {
		//		thread.Trait = callBlockTrait
		//	}
		//}
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
