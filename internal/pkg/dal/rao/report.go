package rao

type ListReportsReq struct {
	TeamID int64 `form:"team_id" binding:"required,gt=0"`
	Page   int   `form:"page,default=1"`
	Size   int   `form:"size,default=10"`

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
	TeamID      int64  `json:"team_id"`
	Name        string `json:"name"`
	TaskMode    int32  `json:"task_mode"`
	TaskType    int32  `json:"task_type"`
	Status      int32  `json:"status"`
	RunTimeSec  int64  `json:"run_time_sec"`
	LastTimeSec int64  `json:"last_time_sec"`
	RunUserID   int64  `json:"run_user_id"`
	RunUserName string `json:"run_user_name"`

	SceneID   int64  `json:"scene_id"`
	SceneName string `json:"scene_name"`
}

type DeleteReportReq struct {
	TeamID   int64 `json:"team_id"`
	ReportID int64 `json:"report_id"`
}

type DeleteReportResp struct {
}

type ListMachineReq struct {
	ReportID int64 `form:"report_id" binding:"required,gt=0"`
}

type ListMachineResp struct {
	StartTimeSec int64     `json:"start_time_sec"`
	EndTimeSec   int64     `json:"end_time_sec"`
	ReportStatus int32     `json:"report_status"`
	Metrics      []*Metric `json:"metrics"`
}

type Metric struct {
	CPU    [][]interface{} `json:"cpu"`
	Mem    [][]interface{} `json:"mem"`
	NetIO  [][]interface{} `json:"net_io"`
	DiskIO [][]interface{} `json:"disk_io"`
}
