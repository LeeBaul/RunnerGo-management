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
	dal.MustInitGRPC()
	dal.MustInitRedis()
	dal.MustInitBigCache()
}
