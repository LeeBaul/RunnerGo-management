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

	r := gin.New()
	router.RegisterRouter(r)

	if err := r.Run(fmt.Sprintf(":%d", conf.Conf.Http.Port)); err != nil {
		panic(err)
	}
}
