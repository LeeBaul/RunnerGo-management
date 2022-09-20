package mao

type ReportTask struct {
	UserID   int64     `json:"user_id" bson:"user_id"`
	UserName string    `json:"user_name" bson:"user_name"`
	ReportID int64     `json:"report_id" bson:"report_id"`
	TaskType int32     `json:"task_type" bson:"task_type"`
	TaskMode int32     `json:"task_mode" bson:"task_mode"`
	ModeConf *ModeConf `json:"mode_conf" bson:"mode_conf"`
}
