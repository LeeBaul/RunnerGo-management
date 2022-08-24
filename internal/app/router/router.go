package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"kp-management/internal/app/middleware"
	"kp-management/internal/pkg/handler"
)

func RegisterRouter(r *gin.Engine) {
	// cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "DELETE", "PUT", "PATCH"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Upgrade", "Origin", "Connection", "Accept-Encoding", "Accept-Language", "Host", "x-requested-with"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// routers
	api := r.Group("/management/api")

	// 用户鉴权
	auth := api.Group("/v1/auth")
	auth.POST("/signup", handler.AuthSignup)
	auth.POST("/login", handler.AuthLogin)
	auth.POST("/refresh_token", handler.AuthRefresh)

	auth.GET("/get_user_settings", handler.GetUserSettings)
	auth.POST("/set_user_settings", handler.SetUserSettings)

	auth.POST("/send_email_verify", handler.AuthSendMailVerify)
	auth.POST("/reset_password", handler.AuthResetPassword)

	// 开启接口鉴权
	api.Use(middleware.JWT())

	// 团队
	team := api.Group("/v1/team")
	team.GET("/list", handler.ListTeam)
	team.GET("/members", handler.TeamMembers)
	team.POST("/invite", handler.InviteMember)
	team.POST("/remove", handler.RemoveMember)

	// 首页
	dashboard := api.Group("/v1/dashboard")
	dashboard.GET("/default", handler.DashboardDefault)
	dashboard.GET("/underway_plans", handler.ListUnderwayPlan)

	// 文件夹
	folder := api.Group("/v1/folder")
	folder.POST("/save", handler.SaveFolder)

	// 接口
	target := api.Group("/v1/target")
	target.POST("/save", handler.SaveTarget)
	target.GET("/list", handler.ListTarget)

	target.GET("/trash_list")
	target.POST("/trash", handler.TrashTarget)
	target.POST("/delete", handler.DeleteTarget)

	// 分组
	group := api.Group("/v1/group")
	group.POST("/save", handler.SaveGroup)

	// 场景
	scene := api.Group("/v1/scene")
	scene.POST("/save")
	scene.GET("/list")

	// 测试计划
	plan := api.Group("/v1/plan")
	plan.GET("/list")

	// 测试报告
	report := api.Group("/v1/report")
	report.GET("/list", handler.ListReports)

	// 操作日志
	operation := api.Group("/v1/operation")
	operation.GET("/list", handler.ListOperation)

	// 机器管理

}
