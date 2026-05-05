package model

import (
	"time"

	"gorm.io/gorm"
)

type Navigation struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Label     string         `gorm:"size:255;not null" json:"label"`
	Link      *string        `gorm:"size:512;default:null" json:"link"`
	LinkType  string         `gorm:"size:32;not null;default:'custom'" json:"link_type"`
	ProjectID *uint64        `gorm:"index;default:null" json:"project_id"`
	PageID    *uint64        `gorm:"index;default:null" json:"page_id"`
	ParentID  *uint64        `gorm:"index;default:null" json:"parent_id"`
	SortOrder int            `gorm:"not null;default:0" json:"sort_order"`
	Status          bool           `gorm:"not null;default:1" json:"status"`
	DisplayPosition string         `gorm:"size:16;not null;default:'header'" json:"display_position"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`

	Children []Navigation `gorm:"foreignKey:ParentID" json:"children"`
	Project  *Project     `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Page     *Page        `gorm:"foreignKey:PageID" json:"page,omitempty"`
}

func (Navigation) TableName() string { return "navigations" }
