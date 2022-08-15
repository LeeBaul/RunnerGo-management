package rao

type DashboardDefaultReq struct {
}

type DashboardDefaultResp struct {
	PlanNum   int64 `json:"plan_num"`
	SceneNum  int64 `json:"scene_num"`
	ReportNum int64 `json:"report_num"`
	APINum    int64 `json:"api_num"`
}
