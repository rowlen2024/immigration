package model

import (
	"time"

	"gorm.io/gorm"
)

type Case struct {
	ID               uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID        *uint64        `gorm:"index" json:"project_id"`
	Project          *Project       `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Name             string         `gorm:"size:128;not null;default:''" json:"name"`
	CountryFrom      string         `gorm:"column:country_from;size:64;not null;default:''" json:"country_from"`
	InvestmentAmount string         `gorm:"size:64;not null;default:''" json:"investment_amount"`
	InvestmentValue  *float64       `gorm:"type:decimal(12,2)" json:"investment_value"`
	ProcessingPeriod string         `gorm:"size:64;not null;default:''" json:"processing_period"`
	Description      string         `gorm:"type:text" json:"description"`
	PhotoURL         string         `gorm:"size:512;not null;default:''" json:"photo_url"`
	SortOrder        int            `gorm:"not null;default:0" json:"sort_order"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

func (Case) TableName() string { return "cases" }
