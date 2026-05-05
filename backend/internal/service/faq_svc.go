package service

import (
	"errors"
	"fmt"
	"time"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

// FAQService handles business logic for FAQ entries.
type FAQService struct {
	repo repository.FAQRepository
}

// NewFAQService creates a new FAQService.
func NewFAQService(repo repository.FAQRepository) *FAQService {
	return &FAQService{repo: repo}
}

// List returns FAQs, optionally filtered by project or global flag.
func (s *FAQService) List(projectID *uint64, isGlobal *bool) ([]dto.FAQResponse, error) {
	results, _, err := s.repo.FindAll(repository.FAQQueryParams{
		ProjectID: projectID,
		IsGlobal:  isGlobal,
		Page:      1,
		PerPage:   1000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list faqs: %w", err)
	}
	return toFAQResponses(results), nil
}

// AdminList returns paginated FAQs with optional project filter and search.
func (s *FAQService) AdminList(projectID *uint64, search string, page, perPage int) ([]dto.FAQResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	results, total, err := s.repo.FindAll(repository.FAQQueryParams{
		ProjectID: projectID,
		Search:    search,
		Page:      page,
		PerPage:   perPage,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list faqs: %w", err)
	}
	return toFAQResponses(results), total, nil
}

// toFAQResponses converts FAQWithProject rows to DTO responses.
func toFAQResponses(rows []repository.FAQWithProject) []dto.FAQResponse {
	result := make([]dto.FAQResponse, len(rows))
	for i, r := range rows {
		result[i] = dto.FAQResponse{
			ID:          r.ID,
			Question:    r.Question,
			Answer:      r.Answer,
			ProjectID:   r.ProjectID,
			ProjectName: r.ProjectName,
			ProjectSlug: r.ProjectSlug,
			IsGlobal:    r.IsGlobal,
			SortOrder:   r.SortOrder,
			CreatedAt:   r.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   r.UpdatedAt.Format(time.RFC3339),
		}
	}
	return result
}

// Create creates a new FAQ entry.
func (s *FAQService) Create(faq *model.FAQ) (*model.FAQ, error) {
	if faq == nil {
		return nil, errors.New("faq is nil")
	}
	if faq.Question == "" {
		return nil, errors.New("faq question is required")
	}
	if faq.Answer == "" {
		return nil, errors.New("faq answer is required")
	}
	if err := s.repo.Create(faq); err != nil {
		return nil, fmt.Errorf("failed to create faq: %w", err)
	}
	return faq, nil
}

// Update updates an existing FAQ entry.
func (s *FAQService) Update(id uint64, faq *model.FAQ) (*model.FAQ, error) {
	if faq == nil {
		return nil, errors.New("faq is nil")
	}
	if id == 0 {
		return nil, errors.New("faq id is required")
	}
	faq.ID = id
	if err := s.repo.Update(faq); err != nil {
		return nil, fmt.Errorf("failed to update faq: %w", err)
	}
	return faq, nil
}

// Delete removes an FAQ entry by ID.
func (s *FAQService) Delete(id uint64) error {
	if id == 0 {
		return errors.New("faq id is required")
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete faq: %w", err)
	}
	return nil
}
