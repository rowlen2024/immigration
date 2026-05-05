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
	repo repository.LeadRepository
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

// AdminList returns paginated leads, optionally filtered by status.
func (s *LeadService) AdminList(page, perPage int, status string) ([]model.Lead, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	leads, total, err := s.repo.FindAll(page, perPage, status)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list leads: %w", err)
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
