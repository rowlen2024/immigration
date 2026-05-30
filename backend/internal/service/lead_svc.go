package service

import (
	"errors"
	"fmt"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"

	"github.com/microcosm-cc/bluemonday"
)

// LeadService handles business logic for customer leads.
type LeadService struct {
	repo        repository.LeadRepository
	projectRepo repository.ProjectRepository
}

// NewLeadService creates a new LeadService with the given repository.
func NewLeadService(repo repository.LeadRepository) *LeadService {
	return &LeadService{repo: repo}
}

var leadSanitizer = bluemonday.StrictPolicy()

// Create creates a new lead from a DTO request, sanitizing text fields.
func (s *LeadService) Create(req *dto.LeadRequest) (*model.Lead, error) {
	if req == nil {
		return nil, errors.New("lead is nil")
	}
	if req.Name == "" {
		return nil, errors.New("lead name is required")
	}
	if req.Phone == "" {
		return nil, errors.New("lead phone is required")
	}

	lead := &model.Lead{
		Name:              leadSanitizer.Sanitize(req.Name),
		Phone:             req.Phone,
		Email:             req.Email,
		InterestedProject: req.InterestedProject,
		Message:           leadSanitizer.Sanitize(req.Message),
	}
	if err := s.repo.Create(lead); err != nil {
		return nil, fmt.Errorf("failed to create lead: %w", err)
	}
	return lead, nil
}

// List returns leads with optional filters and pagination.
func (s *LeadService) List(req dto.LeadListRequest) ([]model.Lead, int64, error) {
	leads, total, err := s.repo.FindAll(repository.LeadFilter{
		Status:            req.Status,
		Name:              req.Name,
		Email:             req.Email,
		InterestedProject: req.InterestedProject,
		Page:              req.Page,
		PerPage:           req.PerPage,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list leads: %w", err)
	}

	// preload 项目名称
	if s.projectRepo != nil && len(leads) > 0 {
		slugSet := make(map[string]struct{})
		for _, l := range leads {
			if l.InterestedProject != "" {
				slugSet[l.InterestedProject] = struct{}{}
			}
		}
		if len(slugSet) > 0 {
			slugs := make([]string, 0, len(slugSet))
			for slug := range slugSet {
				slugs = append(slugs, slug)
			}
			projects, err := s.projectRepo.FindBySlugsLight(slugs)
			if err == nil {
				nameBySlug := make(map[string]string, len(projects))
				for _, p := range projects {
					nameBySlug[p.Slug] = p.Name
				}
				for i := range leads {
					if name, ok := nameBySlug[leads[i].InterestedProject]; ok {
						leads[i].ProjectName = name
					}
				}
			}
		}
	}

	return leads, total, nil
}

// Update updates the status and notes of a lead by ID.
func (s *LeadService) Update(id uint64, status, notes string) (*model.Lead, error) {
	if id == 0 {
		return nil, errors.New("lead id is required")
	}
	validStatuses := map[string]bool{
		"new": true, "contacted": true, "qualified": true, "closed": true,
	}
	if !validStatuses[status] {
		return nil, fmt.Errorf("invalid status: %s", status)
	}
	if err := s.repo.UpdateStatus(id, status, notes); err != nil {
		return nil, fmt.Errorf("failed to update lead status: %w", err)
	}
	return &model.Lead{ID: id, Status: status, Notes: notes}, nil
}
