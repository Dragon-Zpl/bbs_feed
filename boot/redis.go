package boot

import (
	"bbs_feed/conf"
	goredis "github.com/go-redis/redis"
	"sync"
)

const (
	SESSION = iota
	CACHE
)

var (
	session *goredis.Client
	cache   *goredis.Client
	rdMu    sync.Mutex
)

func ConnectRedis() {
	var option0 = goredis.Options{
		Network:    "tcp",
		Addr:       conf.RedisConf.Host + ":" + conf.RedisConf.Port,
		Password:   conf.RedisConf.Password,
		DB:         conf.RedisConf.DB,
		MaxRetries: 3,
		PoolSize:   32,
	}
	var option1 = option0
	session = goredis.NewClient(&option0)
	option1.DB = 1
	cache = goredis.NewClient(&option1)
}

// 获取redis cli 句柄
func InstanceRedisCli(db int) *goredis.Client {
	var tmp *goredis.Client
	if db == 0 {
		tmp = session
	} else {
		tmp = cache
	}
	if tmp != nil {
		return tmp
	} else {
		rdMu.Lock()
		defer rdMu.Unlock()
		ConnectRedis()
		if db == 0 {
			tmp = session
		} else {
			tmp = cache
		}
		return tmp
	}
}

func KeyExist(key string) bool {
	count, _ := InstanceRedisCli(1).Exists(key).Result()
	return count == 1
}
