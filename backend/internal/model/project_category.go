package model

import "time"

const (
	ProjectCategoryDraft     = "draft"
	ProjectCategoryPublished = "published"
)

type ProjectCategory struct {
	ID             uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string     `json:"name"`
	Slug           string     `json:"slug"`
	Intro          string     `json:"intro"`
	SEOTitle       string     `gorm:"column:seo_title" json:"seo_title"`
	SEODescription string     `gorm:"column:seo_description" json:"seo_description"`
	Status         string     `json:"status"`
	PublishedAt    *time.Time `json:"published_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type ProjectCategoryMembership struct {
	ProjectID  uint64 `gorm:"primaryKey;autoIncrement:false"`
	CategoryID uint64 `gorm:"primaryKey;autoIncrement:false"`
	SortOrder  int32
}
