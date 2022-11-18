package dal

import (
	"github.com/go-omnibus/proof"
	"github.com/go-redis/redis/v8"

	"kp-management/internal/pkg/conf"
)

var rdbReport *redis.Client

func MustInitRedisForReport() {
	rdbReport = redis.NewClient(&redis.Options{
		Addr:     conf.Conf.RedisReport.Address,
		Password: conf.Conf.RedisReport.Password,
		DB:       conf.Conf.RedisReport.DB,
	})

	proof.Infof("当前redis是：", conf.Conf.RedisReport.Address)
}

func GetRDBForReport() *redis.Client {
	return rdbReport
}
