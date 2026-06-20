package service

import (
	"errors"
	"fmt"
	"time"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

// FAQService handles business logic for FAQ entries.
type FAQService struct {
	repo        repository.FAQRepository
	versionRepo *repository.PublicVersionRepo
}

// NewFAQService creates a new FAQService.
func NewFAQService(repo repository.FAQRepository) *FAQService {
	return &FAQService{repo: repo}
}

func (s *FAQService) RegisterPublicVersions(reg *PublicVersionRegistry) {
	resolver := func(string) (repository.PublicVersion, error) {
		return s.versionRepo.VersionFromQuery(`
SELECT MAX(updated_at) AS updated_at, COUNT(*) AS count FROM (
  SELECT faqs.updated_at AS updated_at FROM faqs
  UNION ALL SELECT projects.updated_at FROM projects INNER JOIN faqs ON faqs.project_id = projects.id WHERE projects.deleted_at IS NULL
) AS versions`)
	}
	reg.Register("public:faqs:list", resolver)
	reg.Register("public:faqs:projects", resolver)
}

// List returns FAQs with optional filtering and pagination.
func (s *FAQService) List(req dto.FAQListRequest) ([]dto.FAQResponse, int64, error) {
	results, total, err := s.repo.FindAll(repository.FAQQueryParams{
		ProjectID: req.ProjectID,
		IsGlobal:  req.IsGlobal,
		Search:    req.Search,
		Page:      req.Page,
		PerPage:   req.PerPage,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list faqs: %w", err)
	}
	return toFAQResponses(results), total, nil
}

func toFAQResponses(rows []repository.FAQWithProject) []dto.FAQResponse {
	result := make([]dto.FAQResponse, len(rows))
	for i, r := range rows {
		result[i] = dto.FAQResponse{
			ID:          r.ID,
			Question:    r.Question,
			Answer:      r.Answer,
			ProjectID:   r.ProjectID,
			ProjectName: r.ProjectName,
			ProjectSlug: r.ProjectSlug,
			IsGlobal:    r.IsGlobal,
			SortOrder:   r.SortOrder,
			CreatedAt:   r.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   r.UpdatedAt.Format(time.RFC3339),
		}
	}
	return result
}

// ListProjects returns projects that have at least one FAQ.
func (s *FAQService) ListProjects() ([]dto.FAQProjectInfo, error) {
	projects, err := s.repo.FindDistinctProjects()
	if err != nil {
		return nil, fmt.Errorf("failed to list faq projects: %w", err)
	}
	result := make([]dto.FAQProjectInfo, len(projects))
	for i, p := range projects {
		result[i] = dto.FAQProjectInfo{
			ID:   p.ID,
			Name: p.Name,
			Slug: p.Slug,
		}
	}
	return result, nil
}

// Create creates a new FAQ entry.
func (s *FAQService) Create(faq *model.FAQ) (*model.FAQ, error) {
	if faq == nil {
		return nil, errors.New("faq is nil")
	}
	if faq.Question == "" {
		return nil, errors.New("faq question is required")
	}
	if faq.Answer == "" {
		return nil, errors.New("faq answer is required")
	}
	if err := s.repo.Create(faq); err != nil {
		return nil, fmt.Errorf("failed to create faq: %w", err)
	}
	return faq, nil
}

// Update updates an existing FAQ entry.
func (s *FAQService) Update(id uint64, req dto.UpdateFAQRequest) (*model.FAQ, error) {
	if id == 0 {
		return nil, errors.New("faq id is required")
	}
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("faq not found: %w", err)
	}
	existing.ProjectID = req.ProjectID
	existing.Question = req.Question
	existing.Answer = req.Answer
	existing.IsGlobal = req.IsGlobal
	existing.SortOrder = req.SortOrder
	if err := s.repo.Update(existing); err != nil {
		return nil, fmt.Errorf("failed to update faq: %w", err)
	}
	return existing, nil
}

// Delete removes an FAQ entry by ID.
func (s *FAQService) Delete(id uint64) error {
	if id == 0 {
		return errors.New("faq id is required")
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete faq: %w", err)
	}
	return nil
}
