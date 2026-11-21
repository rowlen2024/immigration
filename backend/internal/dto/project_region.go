package dto

import "mygo-immigration/backend/internal/model"

type AdminProject struct {
	model.Project
	RegionID    *uint64                 `json:"region_id"`
	CategoryIDs []uint64                `json:"category_ids"`
	Categories  []ProjectCategoryOption `json:"categories"`
}

type CreateProjectRequest struct {
	model.Project
	RegionID    *uint64  `json:"region_id"`
	CategoryIDs []uint64 `json:"category_ids" binding:"omitempty,dive,gt=0"`
	Status      *int8    `json:"status"`
}

type ProjectRegion struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	FlagEmoji string `json:"flag_emoji"`
}

type PublishedProjectRegion struct {
	ProjectID uint64         `json:"project_id"`
	Region    *ProjectRegion `json:"region"`
}

type ProjectRegionOption struct {
	ID uint64 `json:"id"`
	ProjectRegion
}
