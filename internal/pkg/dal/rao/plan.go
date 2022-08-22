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
	PlanID     int64  `json:"plan_id"`
	Name       string `json:"name"`
	UpdatedSec int64  `json:"updated_sec"`
}
