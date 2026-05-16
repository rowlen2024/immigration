package model

import (
	"time"

	"gorm.io/gorm"
)

type Lawyer struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	PhotoURL  string         `gorm:"size:512;not null;default:''" json:"photo_url"`
	Name      string         `gorm:"size:64;not null;default:''" json:"name"`
	Title     string         `gorm:"size:128;not null;default:''" json:"title"`
	Tags      string         `gorm:"size:512;not null;default:''" json:"tags"`
	SortOrder int            `gorm:"not null;default:0" json:"sort_order"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (Lawyer) TableName() string { return "lawyers" }
