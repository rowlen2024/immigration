package model

import (
	"time"

	"gorm.io/gorm"
)

type Page struct {
	ID              uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID       *uint64        `gorm:"index" json:"project_id"`
	Title           string         `gorm:"size:255;not null" json:"title"`
	Slug            string         `gorm:"uniqueIndex;size:255;not null" json:"slug"`
	Content         string         `gorm:"type:longtext" json:"content"`
	CoverImage      string         `gorm:"size:512;not null;default:''" json:"cover_image"`

	// 变体信息（不存数据库，仅 API 输出）
	CoverImageVariants map[string]ImageVariantInfo `gorm:"-" json:"cover_image_variants,omitempty"`
	MetaTitle          string         `gorm:"size:128;not null;default:''" json:"meta_title"`
	MetaDescription string         `gorm:"size:512;not null;default:''" json:"meta_description"`
	Template        string         `gorm:"size:64;not null;default:'default'" json:"template"`
	PageType        string         `gorm:"size:32;not null;default:'default'" json:"page_type"`
	Status          string         `gorm:"type:enum('draft','published');not null;default:'draft'" json:"status"`
	SortOrder       int            `gorm:"not null;default:0" json:"sort_order"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (Page) TableName() string { return "pages" }
