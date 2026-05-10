package service

import (
	"errors"
	"fmt"

	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

// CaseService handles business logic for immigration case studies.
type CaseService struct {
	repo repository.CaseRepository
}

// NewCaseService creates a new CaseService with the given repository.
func NewCaseService(repo repository.CaseRepository) *CaseService {
	return &CaseService{repo: repo}
}

// List returns all cases.
func (s *CaseService) List() ([]model.Case, error) {
	cases, err := s.repo.FindAll("")
	if err != nil {
		return nil, fmt.Errorf("failed to list cases: %w", err)
	}
	return cases, nil
}

// AdminList returns paginated cases.
func (s *CaseService) AdminList(page, perPage int, search string) ([]model.Case, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	cases, err := s.repo.FindAll(search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list cases: %w", err)
	}

	total := int64(len(cases))
	start := (page - 1) * perPage
	if start >= len(cases) {
		return []model.Case{}, total, nil
	}
	end := start + perPage
	if end > len(cases) {
		end = len(cases)
	}
	return cases[start:end], total, nil
}

// Create creates a new case study.
func (s *CaseService) Create(c *model.Case) (*model.Case, error) {
	if c == nil {
		return nil, errors.New("case is nil")
	}
	if c.Name == "" {
		return nil, errors.New("case name is required")
	}
	if err := s.repo.Create(c); err != nil {
		return nil, fmt.Errorf("failed to create case: %w", err)
	}
	return c, nil
}

// Update updates an existing case study.
func (s *CaseService) Update(id uint64, c *model.Case) (*model.Case, error) {
	if c == nil {
		return nil, errors.New("case is nil")
	}
	if id == 0 {
		return nil, errors.New("case id is required")
	}
	c.ID = id
	if err := s.repo.Update(c); err != nil {
		return nil, fmt.Errorf("failed to update case: %w", err)
	}
	return c, nil
}

// Delete removes a case study by ID.
func (s *CaseService) Delete(id uint64) error {
	if id == 0 {
		return errors.New("case id is required")
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete case: %w", err)
	}
	return nil
}
