package redis_ops

import (
	"bbs_feed/boot"
	"encoding/json"
	"github.com/go-redis/redis"
)

// zset add data
func ZAddSort(key string, datas []interface{}) error {
	boot.InstanceRedisCli(boot.CACHE).Del(key)
	var count = len(datas)
	zAddMems := make([]redis.Z, 0, count)
	for i, data := range datas {
		if byteData, err := json.Marshal(data); err == nil {
			zAddMems = append(zAddMems, redis.Z{
				Score:  float64(i),
				Member: string(byteData),
			})
		}

	}
	err := boot.InstanceRedisCli(boot.CACHE).ZAdd(key, zAddMems...).Err()
	return err
}

// sort del data
func DelZAdd(key string, member string) {
	cache := boot.InstanceRedisCli(boot.CACHE)
	rank, err := cache.ZRank(key, member).Result()
	if err != nil {
		return
	}
	if err = cache.ZRem(key, member).Err(); err == nil {
		if zs, err := cache.ZRangeWithScores(key, rank, -1).Result(); err == nil {
			for _, z := range zs {
				cache.ZIncr(key, redis.Z{
					Score:  -1,
					Member: z.Member,
				})
			}
		}
	}
}
