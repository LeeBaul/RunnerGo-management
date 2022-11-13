package consts

const (
	PlanStatusNormal   = 1 // 未开始
	PlanStatusUnderway = 2 // 进行中

	PlanTaskTypeNormal  = 1 // 普通任务
	PlanTaskTypeCronjob = 2 // 定时任务
	PlanTaskTypeMix     = 3 // 混合模式

	PlanModeConcurrence  = 1 // 并发模式
	PlanModeStep         = 2 // 阶梯模式
	PlanModeErrorRate    = 3 // 错误率模式
	PlanModeResponseTime = 4 // 响应时间模式
	PlanModeRPS          = 5 //每秒请求数模式
	PlanModeTPS          = 6 //每秒事务数模式
	PlanModeMix          = 7 // 混合模式

	// 定时任务的几个状态
	TimedTaskWaitEnable = 0 // 未启用
	TimedTaskInExec     = 1 // 运行中
	TimedTaskTimeout    = 2 // 已过期
)
