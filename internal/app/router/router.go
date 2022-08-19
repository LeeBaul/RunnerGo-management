package router

import (
	"kp-management/internal/app/middleware"
	"kp-management/internal/pkg/handler"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	// cors
	r.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"*"},

		AllowMethods:     []string{"POST", "GET", "OPTIONS", "DELETE", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
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

	auth.POST("/send_email_verify", handler.AuthSendMailVerify)
	auth.POST("/reset_password", handler.AuthResetPassword)

	// 开启接口鉴权
	api.Use(middleware.JWT())

	// 团队
	team := api.Group("/v1/team")
	team.GET("/list", handler.ListTeam)
	team.GET("/members", handler.TeamMembers) // 邀请人
	//team.POST("/invite", handler.InviteMember)
	// 移出成员

	// 首页
	dashboard := api.Group("/v1/dashboard")
	dashboard.GET("/default", handler.DashboardDefault)
	//测试报告，运行中

	// 文件夹
	folder := api.Group("/v1/folder")
	folder.POST("/save", handler.SaveFolder)

	// 接口
	target := api.Group("/v1/target")
	target.POST("/save", handler.SaveTarget)
	target.POST("/delete", handler.DeleteTarget)
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
