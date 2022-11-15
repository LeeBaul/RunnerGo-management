package rao

type SavePreinstallReq struct {
	ID       int32  `json:"id"`
	TeamID   int64  `json:"team_id" binding:"required"`
	ConfName string `json:"conf_name" binding:"required"`
	//UserID        int64          `json:"user_id" binding:"required"`
	//UserName      string         `json:"user_name" binding:"required"`
	TaskType      int32          `json:"task_type" binding:"required"`
	TaskMode      int32          `json:"task_mode" binding:"required"`
	ModeConf      *ModeConf      `json:"mode_conf" binding:"required"`
	TimedTaskConf *TimedTaskConf `json:"timed_task_conf"`
}

type GetPreinstallDetailReq struct {
	ID int32 `json:"id"`
}

type PreinstallDetailResponse struct {
	ID            int32          `json:"id"`
	TeamID        int64          `json:"team_id" binding:"required"`
	ConfName      string         `json:"conf_name" binding:"required"`
	UserName      string         `json:"user_name" binding:"required"`
	TaskType      int32          `json:"task_type" binding:"required"`
	TaskMode      int32          `json:"task_mode" binding:"required"`
	ModeConf      *ModeConf      `json:"mode_conf" binding:"required"`
	TimedTaskConf *TimedTaskConf `json:"timed_task_conf"`
}

type GetPreinstallListReq struct {
	TeamID int64 `json:"team_id" binding:"required"`
}
