package service

import (
	"fmt"
	"strings"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

type LawyerService struct {
	repo repository.LawyerRepository
}

func NewLawyerService(repo repository.LawyerRepository) *LawyerService {
	return &LawyerService{repo: repo}
}

func (s *LawyerService) List(req dto.LawyerListRequest) ([]model.Lawyer, int64, error) {
	items, total, err := s.repo.FindAll(repository.LawyerFilter{
		Name:    req.Name,
		Page:    req.Page,
		PerPage: req.PerPage,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list lawyers: %w", err)
	}
	for i := range items {
		items[i].PhotoVariants = ResolveImageVariants(items[i].PhotoURL, UploadContextLawyer)
	}
	return items, total, nil
}

func (s *LawyerService) GetByID(id uint64) (*model.Lawyer, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("lawyer not found: %w", err)
	}
	return item, nil
}

func (s *LawyerService) Create(input *dto.CreateLawyerInput) (*model.Lawyer, error) {
	item := &model.Lawyer{
		PhotoURL:  input.PhotoURL,
		Name:      input.Name,
		Title:     input.Title,
		Tags:      strings.Join(input.Tags, ","),
		SortOrder: input.SortOrder,
	}
	if err := s.repo.Create(item); err != nil {
		return nil, fmt.Errorf("failed to create lawyer: %w", err)
	}
	return item, nil
}

func (s *LawyerService) Update(id uint64, input *dto.CreateLawyerInput) (*model.Lawyer, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("lawyer not found: %w", err)
	}
	item.PhotoURL = input.PhotoURL
	item.Name = input.Name
	item.Title = input.Title
	item.Tags = strings.Join(input.Tags, ",")
	item.SortOrder = input.SortOrder
	if err := s.repo.Update(item); err != nil {
		return nil, fmt.Errorf("failed to update lawyer: %w", err)
	}
	return item, nil
}

func (s *LawyerService) Delete(id uint64) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete lawyer: %w", err)
	}
	return nil
}
