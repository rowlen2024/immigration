package service

import (
	"errors"
	"fmt"
	"strings"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

type TestimonialService struct {
	repo          repository.TestimonialRepository
	homeConfigSvc *HomeConfigService
}

func (s *TestimonialService) List(req dto.TestimonialListRequest) ([]model.Testimonial, int64, error) {
	items, total, err := s.repo.FindAll(repository.TestimonialFilter{
		ProjectID: req.ProjectID,
		Nickname:  req.Nickname,
		Rating:    req.Rating,
		Page:      req.Page,
		PerPage:   req.PerPage,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list testimonials: %w", err)
	}
	return items, total, nil
}

func (s *TestimonialService) Delete(id uint64) error {
	if id == 0 {
		return errors.New("testimonial id is required")
	}
	if s.homeConfigSvc != nil {
		if err := s.homeConfigSvc.RemoveFeaturedTestimonialID(id); err != nil {
			return fmt.Errorf("failed to clean up featured testimonial ref: %w", err)
		}
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete testimonial: %w", err)
	}
	return nil
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

func (s *TestimonialService) Update(id uint64, req dto.UpdateTestimonialRequest) (*model.Testimonial, error) {
	if id == 0 {
		return nil, errors.New("testimonial id is required")
	}
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("testimonial not found: %w", err)
	}
	if strings.TrimSpace(req.Nickname) == "" {
		return nil, errors.New("nickname is required")
	}
	if strings.TrimSpace(req.Content) == "" {
		return nil, errors.New("content is required")
	}
	if req.Rating < 1 || req.Rating > 5 {
		req.Rating = 5
	}
	existing.AvatarURL = req.AvatarURL
	existing.Nickname = req.Nickname
	existing.Rating = req.Rating
	existing.Content = req.Content
	existing.SortOrder = req.SortOrder
	if err := s.repo.Update(existing); err != nil {
		return nil, fmt.Errorf("failed to update testimonial: %w", err)
	}
	return existing, nil
}

