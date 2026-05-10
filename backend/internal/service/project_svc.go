package service

import (
	"errors"
	"fmt"
	"strings"

	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

// ProjectService handles business logic for immigration projects.
type ProjectService struct {
	repo    repository.ProjectRepository
	navRepo repository.NavigationRepository
}

// NewProjectService creates a new ProjectService with the given repository.
func NewProjectService(repo repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

// GetBySlug returns a project by its slug with all relations preloaded.
func (s *ProjectService) GetBySlug(slug string) (*model.Project, error) {
	if slug == "" {
		return nil, errors.New("slug is required")
	}
	project, err := s.repo.FindBySlug(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get project by slug: %w", err)
	}
	return project, nil
}

// List returns paginated projects.
func (s *ProjectService) List(page, perPage int, search, status string) ([]model.Project, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}
	projects, total, err := s.repo.FindAll(page, perPage, search, status)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list projects: %w", err)
	}
	return projects, total, nil
}

// AdminList returns paginated projects (alias for List).
func (s *ProjectService) AdminList(page, perPage int, search, status string) ([]model.Project, int64, error) {
	return s.List(page, perPage, search, status)
}

// Compare returns multiple projects by their slugs for side-by-side comparison.
func (s *ProjectService) Compare(slugs []string) ([]model.Project, error) {
	if len(slugs) == 0 {
		return nil, errors.New("at least one slug is required")
	}
	if len(slugs) > 5 {
		return nil, errors.New("cannot compare more than 5 projects at once")
	}
	projects, err := s.repo.FindBySlugs(slugs)
	if err != nil {
		return nil, fmt.Errorf("failed to compare projects: %w", err)
	}
	return projects, nil
}

// CompareRow represents a single comparison row.
type CompareRow struct {
	Label string `json:"label"`
	A     string `json:"a"`
	B     string `json:"b"`
}

// CompareResult holds the full comparison output.
type CompareResult struct {
	Projects []CompareProject `json:"projects"`
	Rows     []CompareRow     `json:"rows"`
}

// CompareProject holds minimal project info for the comparison header.
type CompareProject struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

// CompareRows returns formatted comparison rows for two projects.
func (s *ProjectService) CompareRows(slugs []string) (*CompareResult, error) {
	projects, err := s.Compare(slugs)
	if err != nil {
		return nil, err
	}
	if len(projects) < 2 {
		return nil, errors.New("需要两个项目进行对比")
	}

	a, b := projects[0], projects[1]

	projInfo := []CompareProject{
		{Title: a.Name, Slug: a.Slug},
		{Title: b.Name, Slug: b.Slug},
	}

	reqsA := joinRequirements(a.Requirements)
	reqsB := joinRequirements(b.Requirements)

	phaseCountA := fmt.Sprintf("%d 个阶段", len(a.TimelinePhases))
	phaseCountB := fmt.Sprintf("%d 个阶段", len(b.TimelinePhases))

	rows := []CompareRow{
		{Label: "投资金额", A: a.InvestmentAmount, B: b.InvestmentAmount},
		{Label: "办理周期", A: a.ProcessingPeriod, B: b.ProcessingPeriod},
		{Label: "适合人群", A: a.TargetCrowd, B: b.TargetCrowd},
		{Label: "申请条件", A: reqsA, B: reqsB},
		{Label: "费用总计", A: a.CostsTotal, B: b.CostsTotal},
		{Label: "流程步骤", A: phaseCountA, B: phaseCountB},
	}

	return &CompareResult{Projects: projInfo, Rows: rows}, nil
}

func joinRequirements(reqs []model.Requirement) string {
	if len(reqs) == 0 {
		return ""
	}
	labels := make([]string, len(reqs))
	for i, r := range reqs {
		prefix := ""
		if r.IsRequired {
			prefix = "✓ "
		} else {
			prefix = "○ "
		}
		labels[i] = prefix + r.Label
	}
	return strings.Join(labels, "；")
}

// Create creates a new project.
func (s *ProjectService) Create(project *model.Project) (*model.Project, error) {
	if project == nil {
		return nil, errors.New("project is nil")
	}
	if project.Slug == "" {
		return nil, errors.New("project slug is required")
	}
	if project.Name == "" {
		return nil, errors.New("project name is required")
	}
	if err := s.repo.Create(project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}
	return project, nil
}

// Update updates an existing project.
func (s *ProjectService) Update(id uint64, project *model.Project) (*model.Project, error) {
	if project == nil {
		return nil, errors.New("project is nil")
	}
	if id == 0 {
		return nil, errors.New("project id is required")
	}

	// Check slug uniqueness (independent of nav references)
	existing, err := s.repo.FindBySlug(project.Slug)
	if err == nil && existing != nil && existing.ID != id {
		return nil, fmt.Errorf("slug %s is already in use", project.Slug)
	}

	if s.navRepo != nil {
		count, err := s.navRepo.CountByProjectID(id)
		if err != nil {
			return nil, fmt.Errorf("failed to check navigation references: %w", err)
		}
		if count > 0 && (existing == nil || existing.ID != id) {
			return nil, fmt.Errorf("%d 个导航项引用了此项目，请先解除引用", count)
		}
	}

	project.ID = id
	if err := s.repo.Update(project); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}
	return project, nil
}

// Delete performs a soft delete on a project by ID.
func (s *ProjectService) Delete(id uint64) error {
	if id == 0 {
		return errors.New("project id is required")
	}
	if s.navRepo != nil {
		count, err := s.navRepo.CountByProjectID(id)
		if err != nil {
			return fmt.Errorf("failed to check navigation references: %w", err)
		}
		if count > 0 {
			return fmt.Errorf("%d 个导航项引用了此项目，请先解除引用", count)
		}
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	return nil
}
