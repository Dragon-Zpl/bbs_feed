package redis_ops

import (
	"bbs_feed/boot"
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

// zset add data
func ZAddSort(key string, datas []interface{}) error {
	boot.InstanceRedisCli(boot.CACHE).Del(key)
	var count = len(datas)
	zAddMems := make([]redis.Z, 0, count)
	for i, data := range datas {
		if byteData, err := json.Marshal(data); err == nil {
			zAddMems = append(zAddMems, redis.Z{
				Score:  float64(count - i),
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

func KeyExist(key string) bool {
	count, _ := boot.InstanceRedisCli(boot.CACHE).Exists(key).Result()
	return count == 1
}

func HSet(key string, field string, val string, expiration time.Duration) {
	if len(field) == 0 || len(val) == 0 {
		return
	}
	var r = boot.InstanceRedisCli(boot.CACHE)
	r.HSet(key, field, val).Result()
	if expiration != -1 {
		r.Expire(key, expiration).Result()
	}
}

func Hdel(key string, fields ...string) {
	if len(fields) == 0 {
		return
	}
	var r = boot.InstanceRedisCli(boot.CACHE)
	for _, field := range fields {
		r.HDel(key, field)
	}
}

func Del(key ...string) (err error) {
	var r = boot.InstanceRedisCli(boot.CACHE)
	_, err = r.Del(key...).Result()
	return
}

//游标查询
func Scan(matchKey string) (keys []string, err error) {
	var (
		ks         []string
		cursor     uint64
		nextCursor uint64 = 0
		r                 = boot.InstanceRedisCli(boot.CACHE)
	)
	for {
		ks, cursor, err = r.Scan(nextCursor, matchKey, 10).Result()
		keys = append(keys, ks...)
		if cursor == 0 {
			break
		}
		nextCursor = cursor
	}
	return
}

func SaveString(key string, value interface{}) error {
	boot.InstanceRedisCli(boot.CACHE).Del(key)
	_, err := boot.InstanceRedisCli(boot.CACHE).Set(key, value, time.Duration(240)*time.Second).Result()
	if err != nil {
		return err
	}
	return nil
}

func GetString(key string) (string, error) {
	result, err := boot.InstanceRedisCli(boot.CACHE).Get(key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}
