package rao

type ListUnderwayPlanReq struct {
	TeamID int64 `form:"team_id"`
	Page   int   `form:"page"`
	Size   int   `form:"size"`
}

type ListUnderwayPlanResp struct {
	Plans []*Plan `json:"plans"`
	Total int64   `json:"total"`
}

type Plan struct {
	PlanID         int64  `json:"plan_id"`
	TeamID         int64  `json:"team_id"`
	Name           string `json:"name"`
	TaskType       int32  `json:"task_type"`
	Mode           int32  `json:"mode"`
	Status         int32  `json:"status"`
	RunUserID      int64  `json:"run_user_id"`
	RunUserName    string `json:"run_user_name"`
	Remark         string `json:"remark"`
	CreatedTimeSec int64  `json:"created_time_sec"`
	UpdatedTimeSec int64  `json:"updated_time_sec"`
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
