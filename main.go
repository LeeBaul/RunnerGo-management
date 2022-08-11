package main

import (
	"fmt"

	"kp-management/internal"
	"kp-management/internal/app/router"
	"kp-management/internal/pkg/conf"

	"github.com/gin-gonic/gin"
)

func main() {
	internal.InitProjects()

	r := gin.Default()
	router.RegisterRouter(r)

	if err := r.Run(fmt.Sprintf(":%d", conf.Conf.Http.Port)); err != nil {
		panic(err)
	}
}
