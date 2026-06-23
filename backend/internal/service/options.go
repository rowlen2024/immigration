package service

import (
	"fmt"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/repository"
)

const defaultOptionsPerPage = 500
const maxOptionsPerPage = 500

func normalizeOptionsPagination(page, perPage int) (int, int) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = defaultOptionsPerPage
	}
	if perPage > maxOptionsPerPage {
		perPage = maxOptionsPerPage
	}
	return page, perPage
}

func (s *ProjectService) Options(req dto.ProjectOptionRequest, publishedOnly bool) ([]dto.ProjectOption, int64, error) {
	page, perPage := normalizeOptionsPagination(req.Page, req.PerPage)
	status := ""
	if publishedOnly {
		status = "1"
	}

	rows, total, err := s.repo.FindOptions(repository.ProjectFilter{
		Name:    req.Name,
		Status:  status,
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list project options: %w", err)
	}

	options := make([]dto.ProjectOption, len(rows))
	for i, row := range rows {
		options[i] = dto.ProjectOption{ID: row.ID, Slug: row.Slug, Name: row.Name}
	}
	return options, total, nil
}

func (s *CaseService) Options(req dto.CaseOptionRequest) ([]dto.CaseOption, int64, error) {
	page, perPage := normalizeOptionsPagination(req.Page, req.PerPage)
	rows, total, err := s.repo.FindOptions(repository.CaseFilter{
		ProjectID: req.ProjectID,
		Name:      req.Name,
		Page:      page,
		PerPage:   perPage,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list case options: %w", err)
	}

	options := make([]dto.CaseOption, len(rows))
	for i, row := range rows {
		options[i] = dto.CaseOption{ID: row.ID, Name: row.Name}
	}
	return options, total, nil
}

func (s *TestimonialService) Options(req dto.TestimonialOptionRequest) ([]dto.TestimonialOption, int64, error) {
	page, perPage := normalizeOptionsPagination(req.Page, req.PerPage)
	rows, total, err := s.repo.FindOptions(repository.TestimonialFilter{
		Nickname: req.Nickname,
		Page:     page,
		PerPage:  perPage,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list testimonial options: %w", err)
	}

	options := make([]dto.TestimonialOption, len(rows))
	for i, row := range rows {
		options[i] = dto.TestimonialOption{ID: row.ID, Nickname: row.Nickname}
	}
	return options, total, nil
}

func (s *PageService) Options(req dto.PageOptionRequest) ([]dto.PageOption, int64, error) {
	page, perPage := normalizeOptionsPagination(req.Page, req.PerPage)
	rows, total, err := s.repo.FindOptions(repository.PageFilter{
		PageType: req.PageType,
		Title:    req.Title,
		Status:   req.Status,
		Page:     page,
		PerPage:  perPage,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list page options: %w", err)
	}

	options := make([]dto.PageOption, len(rows))
	for i, row := range rows {
		options[i] = dto.PageOption{ID: row.ID, Slug: row.Slug, Title: row.Title}
	}
	return options, total, nil
}
