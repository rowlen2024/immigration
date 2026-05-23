package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Media struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Filename     string         `gorm:"size:255;not null" json:"filename"`
	OriginalName string         `gorm:"size:255;not null;default:''" json:"original_name"`
	URL          string         `gorm:"size:512;not null;default:''" json:"url"`
	MimeType     string         `gorm:"size:64;not null;default:''" json:"mime_type"`
	SizeBytes    uint64         `gorm:"not null;default:0" json:"size_bytes"`
	Variants     datatypes.JSON `gorm:"type:json" json:"variants"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt    time.Time      `json:"created_at"`
}

func (Media) TableName() string { return "media" }
