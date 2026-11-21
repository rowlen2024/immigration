package service

import (
	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
)

type projectCategoryWriter interface {
	CreateWithCategories(project *model.Project, categoryIDs []uint64, requestedStatus *int8) error
	UpdateWithCategories(project *model.Project, categoryIDs []uint64) error
}

func (s *ProjectService) AdminRows(projects []model.Project) ([]dto.AdminProject, error) {
	rows := make([]dto.AdminProject, len(projects))
	ids := make([]uint64, len(projects))
	rowIndex := make(map[uint64]int, len(projects))
	for i, project := range projects {
		rows[i] = dto.AdminProject{
			Project: project, RegionID: project.RegionID,
			CategoryIDs: []uint64{}, Categories: []dto.ProjectCategoryOption{},
		}
		ids[i] = project.ID
		rowIndex[project.ID] = i
	}
	if s.categoryRepo == nil || len(ids) == 0 {
		return rows, nil
	}

	assignments, err := s.categoryRepo.AssignmentsByProjectIDs(ids)
	if err != nil {
		return nil, err
	}
	for _, item := range assignments {
		row := &rows[rowIndex[item.ProjectID]]
		row.CategoryIDs = append(row.CategoryIDs, item.ID)
		row.Categories = append(row.Categories, dto.ProjectCategoryOption{
			ID: item.ID, Name: item.Name, Slug: item.Slug, Status: item.Status,
		})
	}
	return rows, nil
}
