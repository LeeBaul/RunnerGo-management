package rao

type ListReportsReq struct {
	TeamID int64 `form:"team_id"`
	Page   int   `form:"page"`
	Size   int   `form:"size"`

	Keyword      string `form:"keyword"`
	StartTimeSec int64  `form:"start_time_sec"`
	EndTimeSec   int64  `form:"end_time_sec"`
}

type ListReportsResp struct {
	Reports []*Report `json:"reports"`
	Total   int64     `json:"total"`
}

type Report struct {
	ReportID    int64  `json:"report_id"`
	Name        string `json:"name"`
	Mode        int32  `json:"mode"`
	Status      int32  `json:"status"`
	RunTimeSec  int64  `json:"run_time_sec"`
	LastTimeSec int64  `json:"last_time_sec"`
	RunUserID   int64  `json:"run_user_id"`
	RunUserName string `json:"run_user_name"`
	TeamID      int64  `json:"team_id"`
	TaskType    int32  `json:"task_type"`
	SceneType   int32  `json:"scene_type"`
}
