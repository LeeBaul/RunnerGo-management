// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameAPI = "api"

// API mapped from table <api>
type API struct {
	ID          int64          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"` // ID
	TargetID    int64          `gorm:"column:target_id;not null" json:"target_id"`        // 目标ID
	URL         string         `gorm:"column:url;not null" json:"url"`
	Header      string         `gorm:"column:header;not null" json:"header"`
	Query       string         `gorm:"column:query;not null" json:"query"`
	Body        string         `gorm:"column:body;not null" json:"body"`
	Auth        string         `gorm:"column:auth;not null" json:"auth"`
	Description string         `gorm:"column:description;not null" json:"description"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName API's table name
func (*API) TableName() string {
	return TableNameAPI
}
