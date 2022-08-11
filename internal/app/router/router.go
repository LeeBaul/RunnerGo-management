package router

import (
	"time"

	"kp-management/internal/pkg/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	// cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "DELETE", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// routers
	api := r.Group("/management/api")

	// 首页
	dashboard := api.Group("/v1/dashboard")
	dashboard.GET("/default", handler.DashboardDefault)

	// 团队
	team := api.Group("/v1/team")
	team.GET("/list", handler.ListTeam)

	// 接口

	// 场景

	// 测试计划
	plan := api.Group("/v1/plan")
	plan.GET("/list/underway", handler.ListUnderwayPlan)

	// 操作日志
	operation := api.Group("/v1/operation")
	operation.GET("/list")

	// 测试报告
	report := api.Group("/v1/report")
	report.GET("/list")

	// 机器管理

}
