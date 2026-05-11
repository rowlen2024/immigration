package service

import (
	"fmt"

	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

// CompareConfigService handles business logic for project comparison configuration.
type CompareConfigService struct {
	repo repository.CompareConfigRepository
}

// GetByProjectID returns the compare config for a project, or nil if not found.
func (s *CompareConfigService) GetByProjectID(projectID uint64) (*model.CompareConfig, error) {
	cfg, err := s.repo.FindByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get compare config: %w", err)
	}
	return cfg, nil
}

// Save creates or updates a compare config.
func (s *CompareConfigService) Save(cfg *model.CompareConfig) (*model.CompareConfig, error) {
	if err := s.repo.Upsert(cfg); err != nil {
		return nil, fmt.Errorf("failed to save compare config: %w", err)
	}
	return cfg, nil
}

// DeleteByProjectID removes the compare config for a project.
func (s *CompareConfigService) DeleteByProjectID(projectID uint64) error {
	if err := s.repo.DeleteByProjectID(projectID); err != nil {
		return fmt.Errorf("failed to delete compare config: %w", err)
	}
	return nil
}
