package dal

import (
	"time"

	"github.com/go-redis/redis"
)

var (
	RDB          *redis.Client
	timeDuration = 3 * time.Second
)

type RedisClient struct {
	Client *redis.Client
}

func InitRedisClient(addr, password string, db int64) (err error) {
	RDB = redis.NewClient(
		&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       int(db),
		})
	_, err = RDB.Ping().Result()
	return err
}

func InsertStatus(key, value string, expiration time.Duration) (err error) {
	if expiration < 20*time.Second {
		expiration = 20 * time.Second
	}
	err = RDB.Set(key, value, expiration).Err()
	if err != nil {
		return
	}
	return
}

func QueryPlanStatus(key string) (err error, value string) {
	value, err = RDB.Get(key).Result()
	if err != nil {
		return
	}
	return
}

// QueryTimingTaskStatus 查询定时任务状态
func QueryTimingTaskStatus(key string) bool {
	ticker := time.NewTicker(timeDuration)
	for {
		select {
		case <-ticker.C:
			value, _ := RDB.Get(key).Result()
			if value == "false" {
				return false
			}
		}
		time.Sleep(timeDuration)
	}
}

func QuerySceneStatus(key string) (err error, value string) {
	value, err = RDB.Get(key).Result()
	if err != nil {
		return
	}
	return
}
