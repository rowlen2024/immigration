package model

import "time"

type ProjectNews struct {
	ProjectID uint64    `gorm:"primaryKey" json:"project_id"`
	PageID    uint64    `gorm:"primaryKey" json:"page_id"`
	CreatedAt time.Time `json:"created_at"`
	Page      *Page     `gorm:"foreignKey:PageID" json:"page,omitempty"`
}

func (ProjectNews) TableName() string { return "project_news" }
