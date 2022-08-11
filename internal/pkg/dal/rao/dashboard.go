package rao

type DashboardDefaultReq struct {
}

type DashboardDefaultResp struct {
	PlanNum   int `json:"plan_num"`
	SceneNum  int `json:"scene_num"`
	ReportNum int `json:"report_num"`
	APINum    int `json:"api_num"`
}
