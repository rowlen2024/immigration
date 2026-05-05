package model

import (
	"time"

	"gorm.io/gorm"
)

type FAQ struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID *uint64        `gorm:"index" json:"project_id"`
	Question  string         `gorm:"size:512;not null" json:"question"`
	Answer    string         `gorm:"type:text;not null" json:"answer"`
	IsGlobal  bool           `gorm:"not null;default:0" json:"is_global"`
	SortOrder int            `gorm:"not null;default:0" json:"sort_order"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (FAQ) TableName() string { return "faqs" }
