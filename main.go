package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"kp-management/internal"
	"kp-management/internal/app/router"
	"kp-management/internal/pkg/conf"
)

func main() {
	internal.InitProjects()

	//resp, err := dal.ClientGRPC().NotifyStopStress(context.TODO(), &services.NotifyStopStressReq{ReportID: 700})
	////resp, err := dal.ClientGRPC().RegisterMachine(context.TODO(), &services.RegisterMachineReq{Region: "bj"})
	//
	//_, _ = resp, err
	//
	//fmt.Sprintf("%+v, %+v", resp, err)

	r := gin.New()
	router.RegisterRouter(r)

	if err := r.Run(fmt.Sprintf(":%d", conf.Conf.Http.Port)); err != nil {
		panic(err)
	}
}
