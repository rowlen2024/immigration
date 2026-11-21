package dto

import "mygo-immigration/backend/internal/model"

type AdminProject struct {
	model.Project
	RegionID *uint64 `json:"region_id"`
}

type CreateProjectRequest struct {
	model.Project
	RegionID *uint64 `json:"region_id"`
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
