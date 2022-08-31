package rao

type DashboardDefaultReq struct {
	TeamID int64 `form:"team_id" binding:"required,gt=0"`
}

type DashboardDefaultResp struct {
	User       *Member      `json:"user"`
	Operations []*Operation `json:"operations"`
	PlanNum    int64        `json:"plan_num"`
	SceneNum   int64        `json:"scene_num"`
	ReportNum  int64        `json:"report_num"`
	APINum     int64        `json:"api_num"`
}
