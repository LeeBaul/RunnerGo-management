package internal

import (
	"kp-management/internal/pkg/biz/proof"
	"kp-management/internal/pkg/conf"
	"kp-management/internal/pkg/dal"
)

func InitProjects() {
	conf.MustInitConf()
	dal.MustInitMySQL()
	dal.MustInitMongo()
	//dal.MustInitElasticSearch()
	proof.MustInitProof()
	//dal.MustInitGRPC()
	dal.MustInitRedis()
	dal.MustInitBigCache()
	
	// 初始化redis客户端
	if err := dal.InitRedisClient(
		conf.Conf.Redis.Address,
		conf.Conf.Redis.Password,
		int64(conf.Conf.Redis.DB),
	); err != nil {
		panic("redis 连接失败")
	}
}
