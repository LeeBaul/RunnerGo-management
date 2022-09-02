package mao

type Task struct {
	PlanID   int64     `bson:"plan_id"`
	ModeConf *ModeConf `bson:"mode_conf"`
}

type ModeConf struct {
	ReheatTime       int64 `bson:"reheat_time"`       // 预热时长
	RoundNum         int64 `bson:"round_num"`         // 轮次
	Concurrency      int64 `bson:"concurrency"`       // 并发数
	ThresholdValue   int64 `bson:"threshold_value"`   // 阈值
	StartConcurrency int64 `bson:"start_concurrency"` // 起始并发数
	Step             int64 `bson:"step"`              // 步长
	StepRunTime      int64 `bson:"step_run_time"`     // 步长执行时长
	MaxConcurrency   int64 `bson:"max_concurrency"`   // 最大并发数
	Duration         int64 `bson:"duration"`          // 稳定持续时长，持续时长
}

type Preinstall struct {
	TeamID   int64     `bson:"team_id"`
	TaskType int32     `bson:"task_type"`
	CronExpr string    `bson:"cron_expr"`
	Mode     int32     `bson:"mode"`
	ModeConf *ModeConf `bson:"mode_conf"`
}
