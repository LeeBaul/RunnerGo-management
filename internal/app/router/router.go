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
	auth.POST("/refresh_token", handler.RefreshToken)

	auth.POST("/send_email_verify", handler.AuthSendMailVerify)
	auth.POST("/reset_password", handler.AuthResetPassword)

	// 开启接口鉴权
	api.Use(middleware.JWT())

	// 用户配置
	setting := api.Group("/v1/setting")
	setting.GET("/get", handler.GetUserSettings)
	setting.POST("/set", handler.SetUserSettings)

	// 团队
	team := api.Group("/v1/team")
	team.GET("/list", handler.ListTeam)
	team.GET("/members", handler.TeamMembers)
	team.POST("/invite", handler.InviteMember)
	team.POST("/remove", handler.RemoveMember)

	// 全局变量
	variable := api.Group("/v1/variable")
	variable.POST("/save", handler.SaveVariable)
	variable.POST("/delete", handler.DeleteVariable)
	variable.POST("/sync", handler.SyncGlobalVariables)
	variable.GET("/list", handler.ListGlobalVariables)
	// 场景变量
	variable.POST("/scene/sync", handler.SyncSceneVariables)
	variable.GET("/scene/list", handler.ListSceneVariables)
	variable.POST("/scene/import", handler.ImportSceneVariables)
	variable.GET("/scene/import/list", handler.ListImportSceneVariables)

	// 首页
	dashboard := api.Group("/v1/dashboard")
	dashboard.GET("/default", handler.DashboardDefault)
	dashboard.GET("/underway_plans", handler.ListUnderwayPlan)

	// 文件夹
	folder := api.Group("/v1/folder")
	folder.POST("/save", handler.SaveFolder)
	folder.GET("/detail", handler.GetFolder)

	// 接口
	target := api.Group("/v1/target")
	// 接口调试
	target.POST("/send", handler.SendTarget)
	target.GET("/result", handler.GetSendTargetResult)
	// 接口保存
	target.POST("/save", handler.SaveTarget)
	target.GET("/list", handler.ListFolderAPI)
	target.GET("/detail", handler.BatchGetTarget)
	// 接口回收站
	target.GET("/trash_list", handler.TrashTargetList)
	target.POST("/trash", handler.TrashTarget)
	target.POST("/recall", handler.RecallTarget)
	target.POST("/delete", handler.DeleteTarget)

	// 分组
	group := api.Group("/v1/group")
	group.POST("/save", handler.SaveGroup)
	group.GET("/detail", handler.GetGroup)

	// 场景
	scene := api.Group("/v1/scene")
	// 场景调试
	scene.POST("/send", handler.SendScene)
	scene.POST("/api/send", handler.SendSceneAPI)
	scene.GET("/result", handler.GetSendSceneResult)
	// 场景管理
	scene.POST("/save", handler.SaveScene)
	scene.GET("/list", handler.ListGroupScene)
	scene.GET("/detail", handler.BatchGetScene)
	scene.GET("/flow/get", handler.GetFlow)
	scene.GET("/flow/batch/get", handler.BatchGetFlow)
	scene.POST("/flow/save", handler.SaveFlow)

	// 测试计划
	plan := api.Group("/v1/plan")
	plan.GET("/list", handler.ListPlans)
	plan.POST("/save", handler.SavePlan)
	plan.GET("/detail", handler.GetPlan)
	plan.POST("/task/save", handler.SavePlanTask)
	plan.POST("/delete", handler.DeletePlan)
	// 计划预设配置
	plan.POST("/preinstall/save", handler.SetPreinstall)
	plan.GET("/preinstall/detail", handler.GetPreinstall)

	// 测试报告
	report := api.Group("/v1/report")
	report.GET("/list", handler.ListReports)
	report.GET("/detail")
	report.GET("/machine", handler.ListMachines)
	report.POST("/delete", handler.DeleteReport)

	// 操作日志
	operation := api.Group("/v1/operation")
	operation.GET("/list", handler.ListOperations)

	// 机器管理

}
