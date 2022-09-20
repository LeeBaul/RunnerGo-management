package mao

type ReportTask struct {
	UserID   int64     `bson:"user_id"`
	UserName string    `bson:"user_name"`
	ReportID int64     `bson:"report_id"`
	TaskType int32     `bson:"task_type"`
	TaskMode int32     `bson:"task_mode"`
	ModeConf *ModeConf `bson:"mode_conf"`
}
