package model

import (
	"time"

	"gorm.io/gorm"
)

type ProjectAdvantage struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID   uint64         `gorm:"not null;index" json:"project_id"`
	Icon        string         `gorm:"size:64;not null;default:''" json:"icon"`
	IconType    string         `gorm:"size:32;not null;default:'lucide'" json:"icon_type"`
	Title       string         `gorm:"size:128;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	SortOrder   int            `gorm:"not null;default:0" json:"sort_order"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (ProjectAdvantage) TableName() string { return "project_advantages" }
