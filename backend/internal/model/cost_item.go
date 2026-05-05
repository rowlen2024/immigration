package model

import (
	"time"

	"gorm.io/gorm"
)

type CostItem struct {
	ID             uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID      uint64         `gorm:"not null;index" json:"project_id"`
	Name           string         `gorm:"size:255;not null" json:"name"`
	Amount         string         `gorm:"size:64;not null;default:''" json:"amount"`
	AmountValue    *float64       `gorm:"type:decimal(12,2)" json:"amount_value"`
	AmountCurrency string         `gorm:"size:3;not null;default:'USD'" json:"amount_currency"`
	Note           string         `gorm:"size:512;not null;default:''" json:"note"`
	SortOrder      int            `gorm:"not null;default:0" json:"sort_order"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CostItem) TableName() string { return "cost_items" }
