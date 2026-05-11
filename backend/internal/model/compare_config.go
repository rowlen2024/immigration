package model

import (
	"time"

	"gorm.io/datatypes"
)

type CompareConfig struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID     uint64         `gorm:"uniqueIndex;not null" json:"project_id"`
	CompareWith   datatypes.JSON `gorm:"type:json;not null" json:"compare_with"`
	CompareFields datatypes.JSON `gorm:"type:json;not null" json:"compare_fields"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (CompareConfig) TableName() string { return "compare_configs" }
