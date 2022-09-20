package mao

type ReportTask struct {
	PlanID   int64     `bson:"plan_id"`
	PlanName string    `bson:"plan_name"`
	ReportID int64     `bson:"report_id"`
	TaskType int32     `bson:"task_type"`
	TaskMode int32     `bson:"task_mode"`
	ModeConf *ModeConf `bson:"mode_conf"`
}
