// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameReport = "report"

// Report mapped from table <report>
type Report struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name      string    `gorm:"column:name;not null" json:"name"`
	Mode      int32     `gorm:"column:mode;not null" json:"mode"`
	Status    int32     `gorm:"column:status;not null" json:"status"`
	RanAt     time.Time `gorm:"column:ran_at;not null" json:"ran_at"`
	RunUserID int64     `gorm:"column:run_user_id;not null" json:"run_user_id"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
	TeamID    int64     `gorm:"column:team_id;not null" json:"team_id"`
}

// TableName Report's table name
func (*Report) TableName() string {
	return TableNameReport
}
