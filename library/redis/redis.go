package redis

import (
	"time"

	"bitcoin/conf"
	"github.com/go-redis/redis"
)

var GoRedisClient *redis.Client

func Init() {
	if !conf.Config.Redis.IsUse {
		return
	}
	GoRedisClient = redis.NewClient(&redis.Options{
		Addr:         conf.Config.Redis.Addr,
		Password:     conf.Config.Redis.Password,
		MaxRetries:   5,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		PoolTimeout:  5 * time.Second,
		DB:           conf.Config.Redis.DB,
		PoolSize:     conf.Config.Redis.PoolSize,
		IdleTimeout:  15 * time.Second,
	})

	pong, err := GoRedisClient.Ping().Result()
	if err != nil || pong != "PONG" {
		panic(err)
	}
}

func Lock(key string, timeout time.Duration) (bool, error) {
	return GoRedisClient.SetNX(key, 1, timeout).Result()
}

func UnLock(key string) (int64, error) {
	return GoRedisClient.Del(key).Result()
}
