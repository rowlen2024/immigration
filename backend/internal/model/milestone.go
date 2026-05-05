package model

import (
	"time"

	"gorm.io/gorm"
)

type Milestone struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID   uint64         `gorm:"not null;index" json:"project_id"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	SortOrder   int            `gorm:"not null;default:0" json:"sort_order"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (Milestone) TableName() string { return "milestones" }
