package rao

type Plan struct {
	PlanID         int64     `json:"plan_id"`
	TeamID         int64     `json:"team_id"`
	Name           string    `json:"name"`
	TaskType       int32     `json:"task_type"`
	Mode           int32     `json:"mode"`
	Status         int32     `json:"status"`
	RunUserID      int64     `json:"run_user_id"`
	RunUserName    string    `json:"run_user_name"`
	Remark         string    `json:"remark"`
	CreatedTimeSec int64     `json:"created_time_sec"`
	UpdatedTimeSec int64     `json:"updated_time_sec"`
	ModeConf       *ModeConf `json:"mode_conf"`
}

type ModeConf struct {
	ReheatTime       int64 `json:"reheat_time"`       // 预热时长
	RoundNum         int64 `json:"round_num"`         // 轮次
	Concurrency      int64 `json:"concurrency"`       // 并发数
	ThresholdValue   int64 `json:"threshold_value"`   // 阈值
	StartConcurrency int64 `json:"start_concurrency"` // 起始并发数
	Step             int64 `json:"step"`              // 步长
	StepRunTime      int64 `json:"step_run_time"`     // 步长执行时长
	MaxConcurrency   int64 `json:"max_concurrency"`   // 最大并发数
	Duration         int64 `json:"duration"`          // 稳定持续时长，持续时长
}

type ListUnderwayPlanReq struct {
	TeamID int64 `form:"team_id"`
	Page   int   `form:"page"`
	Size   int   `form:"size"`
}

type ListUnderwayPlanResp struct {
	Plans []*Plan `json:"plans"`
	Total int64   `json:"total"`
}

type ListPlansReq struct {
	TeamID int64 `form:"team_id"`
	Page   int   `form:"page"`
	Size   int   `form:"size"`

	Keyword      string `form:"keyword"`
	StartTimeSec int64  `form:"start_time_sec"`
	EndTimeSec   int64  `form:"end_time_sec"`
}

type ListPlansResp struct {
	Plans []*Plan `json:"plans"`
	Total int64   `json:"total"`
}

type SavePlanReq struct {
	PlanID   int64     `json:"plan_id"`
	TeamID   int64     `json:"team_id"`
	Name     string    `json:"name"`
	TaskType int32     `json:"task_type"`
	Mode     int32     `json:"mode"`
	Remark   string    `json:"remark"`
	ModeConf *ModeConf `json:"mode_conf"`
}

type SavePlanResp struct {
}

type GetPlanReq struct {
	TeamID int64 `form:"team_id"`
	PlanID int64 `form:"plan_id"`
}

type GetPlanResp struct {
	Plan *Plan `json:"plan"`
}
