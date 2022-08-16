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

	// 用户鉴权
	auth := r.Group("/v1/auth")
	auth.POST("/signup", handler.AuthSignup)
	auth.POST("/login", handler.AuthLogin)

	// 团队
	team := api.Group("/v1/team")
	team.GET("/list", handler.ListTeam)
	team.GET("/members", handler.TeamMembers)
	team.POST("/invite", handler.InviteMember)

	// 首页
	dashboard := api.Group("/v1/dashboard")
	dashboard.GET("/default", handler.DashboardDefault)

	// 文件夹
	folder := api.Group("/v1/folder")
	folder.POST("/save", handler.SaveFolder)
	folder.POST("/delete", handler.DeleteFolder)

	// 接口
	target := api.Group("/v1/target")
	target.POST("/save", handler.SaveTarget)
	target.GET("/list", handler.ListTarget)

	// 分组
	group := api.Group("/v1/group")
	group.POST("/save", handler.SaveGroup)

	// 场景
	scene := api.Group("/v1/scene")
	scene.POST("/save")

	// 测试计划
	plan := api.Group("/v1/plan")
	plan.GET("/list/underway", handler.ListUnderwayPlan)

	// 操作日志
	operation := api.Group("/v1/operation")
	operation.GET("/list", handler.ListOperation)

	// 测试报告
	report := api.Group("/v1/report")
	report.GET("/list")

	// 机器管理

}
