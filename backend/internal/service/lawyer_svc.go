package service

import (
	"fmt"
	"strings"

	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

type LawyerService struct {
	repo *repository.LawyerRepo
}

// LawyerResponse is the API-facing representation with tags as array.
type LawyerResponse struct {
	ID        uint64   `json:"id"`
	PhotoURL  string   `json:"photo_url"`
	Name      string   `json:"name"`
	Title     string   `json:"title"`
	Tags      []string `json:"tags"`
	SortOrder int      `json:"sort_order"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

func toLawyerResponse(l *model.Lawyer) LawyerResponse {
	tags := []string{}
	if l.Tags != "" {
		for _, t := range strings.Split(l.Tags, ",") {
			trimmed := strings.TrimSpace(t)
			if trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}
	return LawyerResponse{
		ID:        l.ID,
		PhotoURL:  l.PhotoURL,
		Name:      l.Name,
		Title:     l.Title,
		Tags:      tags,
		SortOrder: l.SortOrder,
		CreatedAt: l.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: l.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (s *LawyerService) List() ([]LawyerResponse, error) {
	items, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list lawyers: %w", err)
	}
	result := make([]LawyerResponse, len(items))
	for i, item := range items {
		result[i] = toLawyerResponse(&item)
	}
	return result, nil
}

func (s *LawyerService) ListPaginated(page, perPage int, search string) ([]LawyerResponse, int64, error) {
	items, total, err := s.repo.FindPaginated(page, perPage, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list lawyers: %w", err)
	}
	result := make([]LawyerResponse, len(items))
	for i, item := range items {
		result[i] = toLawyerResponse(&item)
	}
	return result, total, nil
}

func (s *LawyerService) GetByID(id uint64) (*LawyerResponse, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("lawyer not found: %w", err)
	}
	resp := toLawyerResponse(item)
	return &resp, nil
}

// CreateLawyerInput is used for both creating and updating.
type CreateLawyerInput struct {
	PhotoURL  string   `json:"photo_url"`
	Name      string   `json:"name"`
	Title     string   `json:"title"`
	Tags      []string `json:"tags"`
	SortOrder int      `json:"sort_order"`
}

func (s *LawyerService) Create(input *CreateLawyerInput) (*LawyerResponse, error) {
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
	resp := toLawyerResponse(item)
	return &resp, nil
}

func (s *LawyerService) Update(id uint64, input *CreateLawyerInput) (*LawyerResponse, error) {
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
	resp := toLawyerResponse(item)
	return &resp, nil
}

func (s *LawyerService) Delete(id uint64) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete lawyer: %w", err)
	}
	return nil
}
