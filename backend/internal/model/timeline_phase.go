package model

import (
	"time"

	"gorm.io/gorm"
)

type TimelinePhase struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID   uint64         `gorm:"not null;index" json:"project_id"`
	PhaseNumber int            `gorm:"not null" json:"phase_number"`
	PhaseName   string         `gorm:"size:128;not null;default:''" json:"phase_name"`
	Duration    string         `gorm:"size:64;not null;default:''" json:"duration"`
	Title       string         `gorm:"size:255;not null;default:''" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	SortOrder   int            `gorm:"not null;default:0" json:"sort_order"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (TimelinePhase) TableName() string { return "timeline_phases" }
