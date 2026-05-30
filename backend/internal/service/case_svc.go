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
	return c, nil
}

func (s *CaseService) List() ([]model.Case, error) {
	return s.ListAll("")
}

func (s *CaseService) ListPaginated(page, perPage int) ([]model.Case, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	cases, total, err := s.repo.FindAllPaginated(page, perPage, "")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list cases: %w", err)
	}
	return cases, total, nil
}

func (s *CaseService) ListAll(search string) ([]model.Case, error) {
	cases, err := s.repo.FindAll(search)
	if err != nil {
		return nil, fmt.Errorf("failed to list cases: %w", err)
	}
	return cases, nil
}

func (s *CaseService) ListFilteredPaginated(projectID *uint64, countryFrom string, page, perPage int) ([]model.Case, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	cases, total, err := s.repo.FindFilteredPaginated(projectID, countryFrom, page, perPage)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list filtered cases: %w", err)
	}
	return cases, total, nil
}

func (s *CaseService) ListByProject(projectID uint64) ([]model.Case, error) {
	cases, err := s.repo.FindByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list cases by project: %w", err)
	}
	return cases, nil
}

func (s *CaseService) AdminList(page, perPage int, search string) ([]model.Case, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	cases, total, err := s.repo.FindAllPaginated(page, perPage, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list cases: %w", err)
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

