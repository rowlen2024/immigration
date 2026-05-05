package model

import (
	"time"

	"gorm.io/gorm"
)

type Requirement struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID uint64         `gorm:"not null;index" json:"project_id"`
	Label     string         `gorm:"size:255;not null" json:"label"`
	IsRequired bool          `gorm:"not null;default:1" json:"is_required"`
	SortOrder int            `gorm:"not null;default:0" json:"sort_order"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (Requirement) TableName() string { return "requirements" }
