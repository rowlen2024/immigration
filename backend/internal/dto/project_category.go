package dto

import (
	"time"

	"mygo-immigration/backend/internal/model"
)

type ProjectCategoryListRequest struct {
	PaginationRequest
	Name   string `form:"name"`
	Status string `form:"status" binding:"omitempty,oneof=draft published"`
}

type SaveProjectCategoryRequest struct {
	Name           string `json:"name" binding:"required,max=128"`
	Slug           string `json:"slug" binding:"required,max=128"`
	Intro          string `json:"intro"`
	SEOTitle       string `json:"seo_title" binding:"max=128"`
	SEODescription string `json:"seo_description" binding:"max=512"`
	Status         string `json:"status" binding:"omitempty,oneof=draft published"`
}

type ProjectCategoryOption struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Status string `json:"status"`
}

type ProjectCategoryPublicSummary struct {
	ID           uint64    `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	ProjectCount int64     `json:"project_count"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ProjectCategoryPublicDetail struct {
	ProjectCategoryPublicSummary
	Intro          string                `json:"intro"`
	SEOTitle       string                `json:"seo_title"`
	SEODescription string                `json:"seo_description"`
	Projects       []ProjectCategoryCard `json:"projects"`
}

type ProjectCategoryCard struct {
	ID                 uint64                            `json:"id"`
	Name               string                            `json:"name"`
	Slug               string                            `json:"slug"`
	Country            string                            `json:"country"`
	FlagEmoji          string                            `json:"flag_emoji"`
	Tagline            string                            `json:"tagline"`
	OverviewText       string                            `json:"overview_text"`
	InvestmentAmount   string                            `json:"investment_amount"`
	ProcessingPeriod   string                            `json:"processing_period"`
	TargetCrowd        string                            `json:"target_crowd"`
	CoverImage         string                            `json:"cover_image"`
	CoverImageVariants map[string]model.ImageVariantInfo `json:"cover_image_variants,omitempty"`
}

type UpdateProjectCategoryProjectRequest struct {
	SortOrder *int32 `json:"sort_order" binding:"required"`
}
