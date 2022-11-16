package record

const (
	OperationOperateCreateFolder = 1  // 创建文件夹
	OperationOperateCreateAPI    = 2  // 创建接口
	OperationOperateCreateGroup  = 3  // 创建分组
	OperationOperateCreatePlan   = 4  // 创建计划
	OperationOperateCreateScene  = 5  // 创建场景
	OperationOperateUpdateFolder = 6  // 修改文件夹
	OperationOperateUpdateAPI    = 7  // 修改接口
	OperationOperateUpdateGroup  = 8  // 修改分组
	OperationOperateUpdatePlan   = 9  // 修改计划
	OperationOperateUpdateScene  = 10 // 修改场景
	OperationOperateClonePlan    = 11 // 克隆计划
	OperationOperateDeleteReport = 12 // 删除报告
	OperationOperateDeleteScene  = 13 // 删除场景
	OperationOperateDeletePlan   = 14 // 删除计划
	OperationOperateRunScene     = 15 // 运行场景
	OperationOperateRunPlan      = 16 // 运行计划

	OperationOperateSavePreinstall   = 17 // 新建预设配置
	OperationOperateUpdatePreinstall = 18 // 修改并保存预设配置
	OperationOperateDeletePreinstall = 19 // 删除预设配置

)
