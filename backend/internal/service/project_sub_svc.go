package service

import (
	"errors"
	"fmt"

	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

type RequirementService struct {
	repo repository.RequirementRepository
}

func NewRequirementService(repo repository.RequirementRepository) *RequirementService {
	return &RequirementService{repo: repo}
}

func (s *RequirementService) List(projectID uint64) ([]model.Requirement, error) {
	items, err := s.repo.FindByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list requirements: %w", err)
	}
	return items, nil
}

func (s *RequirementService) Create(projectID uint64, item *model.Requirement) (*model.Requirement, error) {
	if item == nil {
		return nil, errors.New("requirement is nil")
	}
	if item.Label == "" {
		return nil, errors.New("label is required")
	}
	item.ProjectID = projectID
	if err := s.repo.Create(item); err != nil {
		return nil, fmt.Errorf("failed to create requirement: %w", err)
	}
	return item, nil
}

func (s *RequirementService) Update(projectID uint64, id uint64, item *model.Requirement) (*model.Requirement, error) {
	if item == nil {
		return nil, errors.New("requirement is nil")
	}
	if id == 0 {
		return nil, errors.New("id is required")
	}
	item.ID = id
	item.ProjectID = projectID
	if err := s.repo.Update(item); err != nil {
		return nil, fmt.Errorf("failed to update requirement: %w", err)
	}
	return item, nil
}

func (s *RequirementService) Delete(projectID uint64, id uint64) error {
	if id == 0 {
		return errors.New("id is required")
	}
	_ = projectID
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete requirement: %w", err)
	}
	return nil
}

type CostItemService struct {
	repo repository.CostItemRepository
}

func NewCostItemService(repo repository.CostItemRepository) *CostItemService {
	return &CostItemService{repo: repo}
}

func (s *CostItemService) List(projectID uint64) ([]model.CostItem, error) {
	items, err := s.repo.FindByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list cost items: %w", err)
	}
	return items, nil
}

func (s *CostItemService) Create(projectID uint64, item *model.CostItem) (*model.CostItem, error) {
	if item == nil {
		return nil, errors.New("cost item is nil")
	}
	if item.Name == "" {
		return nil, errors.New("name is required")
	}
	item.ProjectID = projectID
	if err := s.repo.Create(item); err != nil {
		return nil, fmt.Errorf("failed to create cost item: %w", err)
	}
	return item, nil
}

func (s *CostItemService) Update(projectID uint64, id uint64, item *model.CostItem) (*model.CostItem, error) {
	if item == nil {
		return nil, errors.New("cost item is nil")
	}
	if id == 0 {
		return nil, errors.New("id is required")
	}
	item.ID = id
	item.ProjectID = projectID
	if err := s.repo.Update(item); err != nil {
		return nil, fmt.Errorf("failed to update cost item: %w", err)
	}
	return item, nil
}

func (s *CostItemService) Delete(projectID uint64, id uint64) error {
	if id == 0 {
		return errors.New("id is required")
	}
	_ = projectID
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete cost item: %w", err)
	}
	return nil
}

type TimelinePhaseService struct {
	repo repository.TimelinePhaseRepository
}

func NewTimelinePhaseService(repo repository.TimelinePhaseRepository) *TimelinePhaseService {
	return &TimelinePhaseService{repo: repo}
}

func (s *TimelinePhaseService) List(projectID uint64) ([]model.TimelinePhase, error) {
	items, err := s.repo.FindByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list timeline phases: %w", err)
	}
	return items, nil
}

func (s *TimelinePhaseService) Create(projectID uint64, item *model.TimelinePhase) (*model.TimelinePhase, error) {
	if item == nil {
		return nil, errors.New("timeline phase is nil")
	}
	if item.Title == "" {
		return nil, errors.New("title is required")
	}
	item.ProjectID = projectID
	if err := s.repo.Create(item); err != nil {
		return nil, fmt.Errorf("failed to create timeline phase: %w", err)
	}
	return item, nil
}

func (s *TimelinePhaseService) Update(projectID uint64, id uint64, item *model.TimelinePhase) (*model.TimelinePhase, error) {
	if item == nil {
		return nil, errors.New("timeline phase is nil")
	}
	if id == 0 {
		return nil, errors.New("id is required")
	}
	item.ID = id
	item.ProjectID = projectID
	if err := s.repo.Update(item); err != nil {
		return nil, fmt.Errorf("failed to update timeline phase: %w", err)
	}
	return item, nil
}

func (s *TimelinePhaseService) Delete(projectID uint64, id uint64) error {
	if id == 0 {
		return errors.New("id is required")
	}
	_ = projectID
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete timeline phase: %w", err)
	}
	return nil
}
