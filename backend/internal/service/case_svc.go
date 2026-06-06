package service

import (
	"errors"
	"fmt"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/logging"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"

	"github.com/google/uuid"
)

type CaseService struct {
	repo          repository.CaseRepository
	homeConfigSvc *HomeConfigService
}

func NewCaseService(repo repository.CaseRepository, homeConfigSvc *HomeConfigService) *CaseService {
	return &CaseService{repo: repo, homeConfigSvc: homeConfigSvc}
}

func (s *CaseService) GetBySlug(slug string) (*model.Case, error) {
	if slug == "" {
		return nil, errors.New("slug is required")
	}
	c, err := s.repo.FindBySlug(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get case by slug: %w", err)
	}
	if c.PhotoURL != "" {
		c.PhotoVariants = ResolveImageVariants(c.PhotoURL, UploadContextCase)
	}
	return c, nil
}

func (s *CaseService) List(req dto.CaseListRequest) ([]model.Case, int64, error) {
	cases, total, err := s.repo.FindAll(repository.CaseFilter{
		ProjectID:   req.ProjectID,
		CountryFrom: req.CountryFrom,
		Name:        req.Name,
		Page:        req.Page,
		PerPage:     req.PerPage,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list cases: %w", err)
	}
	for i := range cases {
		if cases[i].PhotoURL != "" {
			cases[i].PhotoVariants = ResolveImageVariants(cases[i].PhotoURL, UploadContextCase)
		}
	}
	return cases, total, nil
}

func (s *CaseService) Create(c *model.Case) (*model.Case, error) {
	if c == nil {
		return nil, errors.New("case is nil")
	}
	if c.Name == "" {
		return nil, errors.New("case name is required")
	}
	if c.Slug == "" {
		c.Slug = uuid.New().String()
	}
	c.Content = HTMLSanitizer.Sanitize(c.Content)
	if err := s.repo.Create(c); err != nil {
		return nil, fmt.Errorf("failed to create case: %w", err)
	}
	return c, nil
}

func (s *CaseService) Update(id uint64, req dto.UpdateCaseRequest) (*model.Case, error) {
	if id == 0 {
		return nil, errors.New("case id is required")
	}
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("case not found: %w", err)
	}
	existing.Name = req.Name
	existing.CountryFrom = req.CountryFrom
	existing.ProjectID = req.ProjectID
	existing.InvestmentAmount = req.InvestmentAmount
	existing.InvestmentValue = req.InvestmentValue
	existing.ProcessingPeriod = req.ProcessingPeriod
	existing.Content = HTMLSanitizer.Sanitize(req.Content)
	existing.PhotoURL = req.PhotoURL
	existing.SortOrder = req.SortOrder
	if err := s.repo.Update(existing); err != nil {
		return nil, fmt.Errorf("failed to update case: %w", err)
	}
	return existing, nil
}

func (s *CaseService) Delete(id uint64) error {
	if id == 0 {
		return errors.New("case id is required")
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete case: %w", err)
	}
	if s.homeConfigSvc != nil {
		if err := s.homeConfigSvc.RemoveFeaturedCaseID(id); err != nil {
			logging.Logger.Warn("home_config: failed to clean up featured case ref after delete",
				"case_id", id, "error", err)
		}
	}
	return nil
}

