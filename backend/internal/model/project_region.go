package model

import "time"

type ProjectRegion struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Code      string    `gorm:"type:char(2);uniqueIndex:idx_project_regions_code;not null" json:"code"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	FlagEmoji string    `gorm:"size:16;not null;default:''" json:"flag_emoji"`
	SortOrder int       `gorm:"not null;default:0" json:"sort_order"`
	UpdatedAt time.Time `gorm:"type:datetime(6);not null;default:CURRENT_TIMESTAMP(6)" json:"-"`
}

func (ProjectRegion) TableName() string { return "project_regions" }
