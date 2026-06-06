package model

import (
	"time"

	"gorm.io/gorm"
)

type Testimonial struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID *uint64        `gorm:"index" json:"project_id"`
	AvatarURL string         `gorm:"size:512;not null;default:''" json:"avatar_url"`
	Nickname  string         `gorm:"size:64;not null;default:''" json:"nickname"`
	Rating    uint8          `gorm:"not null;default:5" json:"rating"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	SortOrder int            `gorm:"not null;default:0" json:"sort_order"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`

	// 变体信息（不存数据库，仅 API 输出）
	AvatarVariants map[string]ImageVariantInfo `gorm:"-" json:"avatar_variants,omitempty"`
}

func (Testimonial) TableName() string { return "testimonials" }
