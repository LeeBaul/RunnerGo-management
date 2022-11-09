package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kp-management/internal"
	"kp-management/internal/app/router"
	"kp-management/internal/pkg/conf"
	"kp-management/internal/pkg/handler"
)

func main() {
	internal.InitProjects()

	r := gin.New()
	router.RegisterRouter(r)

	// 异步执行定时任务
	go func() {
		handler.TimedTaskExec()
	}()

	if err := r.Run(fmt.Sprintf(":%d", conf.Conf.Http.Port)); err != nil {
		panic(err)
	}
}
