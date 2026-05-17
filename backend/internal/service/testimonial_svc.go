package service

import (
	"errors"
	"fmt"
	"strings"

	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

type TestimonialService struct {
	repo repository.TestimonialRepository
}

func (s *TestimonialService) ListByProject(projectID uint64) ([]model.Testimonial, error) {
	items, err := s.repo.FindByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list testimonials: %w", err)
	}
	return items, nil
}

func (s *TestimonialService) ListAll() ([]model.Testimonial, error) {
	items, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list testimonials: %w", err)
	}
	return items, nil
}

func (s *TestimonialService) Create(projectID uint64, t *model.Testimonial) (*model.Testimonial, error) {
	if t == nil {
		return nil, errors.New("testimonial is nil")
	}
	if strings.TrimSpace(t.Nickname) == "" {
		return nil, errors.New("nickname is required")
	}
	if strings.TrimSpace(t.Content) == "" {
		return nil, errors.New("content is required")
	}
	if t.Rating < 1 || t.Rating > 5 {
		t.Rating = 5
	}
	t.ProjectID = &projectID
	if err := s.repo.Create(t); err != nil {
		return nil, fmt.Errorf("failed to create testimonial: %w", err)
	}
	return t, nil
}

func (s *TestimonialService) Update(id uint64, t *model.Testimonial) (*model.Testimonial, error) {
	if t == nil {
		return nil, errors.New("testimonial is nil")
	}
	if id == 0 {
		return nil, errors.New("testimonial id is required")
	}
	if strings.TrimSpace(t.Nickname) == "" {
		return nil, errors.New("nickname is required")
	}
	if strings.TrimSpace(t.Content) == "" {
		return nil, errors.New("content is required")
	}
	if t.Rating < 1 || t.Rating > 5 {
		t.Rating = 5
	}
	t.ID = id
	if err := s.repo.Update(t); err != nil {
		return nil, fmt.Errorf("failed to update testimonial: %w", err)
	}
	return t, nil
}

func (s *TestimonialService) HardDelete(id uint64) error {
	if id == 0 {
		return errors.New("testimonial id is required")
	}
	if err := s.repo.HardDelete(id); err != nil {
		return fmt.Errorf("failed to delete testimonial: %w", err)
	}
	return nil
}
