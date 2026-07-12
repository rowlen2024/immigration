package model

import (
	"time"

	"gorm.io/gorm"
)

type Case struct {
	ID               uint64   `gorm:"primaryKey;autoIncrement" json:"id"`
	Slug             string   `gorm:"uniqueIndex;size:36;not null" json:"slug"`
	ProjectID        *uint64  `gorm:"index" json:"project_id"`
	Project          *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Name             string   `gorm:"size:128;not null;default:''" json:"name"`
	CountryFrom      string   `gorm:"column:country_from;size:64;not null;default:''" json:"country_from"`
	InvestmentAmount string   `gorm:"size:64;not null;default:''" json:"investment_amount"`
	InvestmentValue  *float64 `gorm:"type:decimal(12,2)" json:"investment_value"`
	ProcessingPeriod string   `gorm:"size:64;not null;default:''" json:"processing_period"`
	Content          string   `gorm:"type:longtext" json:"content"`
	PhotoURL         string   `gorm:"size:512;not null;default:''" json:"photo_url"`
	SortOrder        int      `gorm:"not null;default:0" json:"sort_order"`
	IsPinned         bool     `gorm:"not null;default:false" json:"is_pinned"`

	// 变体信息（不存数据库，仅 API 输出）
	PhotoVariants map[string]ImageVariantInfo `gorm:"-" json:"photo_variants,omitempty"`
	DeletedAt     gorm.DeletedAt              `gorm:"index" json:"-"`
	CreatedAt     time.Time                   `json:"created_at"`
	UpdatedAt     time.Time                   `json:"updated_at"`
}

func (Case) TableName() string { return "cases" }
