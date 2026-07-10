package service

import (
	"errors"
	"fmt"
	"mygo-immigration/backend/internal/config"
	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/logging"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
	"strings"
)

// ProjectService handles business logic for immigration projects.
type ProjectService struct {
	repo              repository.ProjectRepository
	navRepo           repository.NavigationRepository
	homeConfigSvc     *HomeConfigService
	requirementRepo   repository.RequirementRepository
	costItemRepo      repository.CostItemRepository
	timelinePhaseRepo repository.TimelinePhaseRepository
	milestoneRepo     repository.MilestoneRepository
	advantageRepo     repository.ProjectAdvantageRepository
	compareConfigRepo repository.CompareConfigRepository
	caseRepo          repository.CaseRepository
	testimonialRepo   repository.TestimonialRepository
	faqRepo           repository.FAQRepository
	versionRepo       *repository.PublicVersionRepo
}

// NewProjectService creates a new ProjectService with the given dependencies.
func NewProjectService(repo repository.ProjectRepository, navRepo repository.NavigationRepository) *ProjectService {
	return &ProjectService{repo: repo, navRepo: navRepo}
}

func (s *ProjectService) RegisterPublicVersions(reg *PublicVersionRegistry) {
	reg.Register("public:projects:list", func(string) (repository.PublicVersion, error) {
		return tableVersion(s.versionRepo, "projects", "deleted_at IS NULL")
	})
	reg.Register("public:project:", s.publicProjectVersion)
	reg.Register("public:compare:", s.publicCompareVersion)
}

func (s *ProjectService) publicCompareVersion(key string) (repository.PublicVersion, error) {
	raw := publicSlug(key, "public:compare:")
	parts := strings.Split(raw, ":")
	versions := make([]repository.PublicVersion, 0, len(parts))
	for _, slug := range parts {
		if slug == "" {
			continue
		}
		v, err := s.publicProjectVersion("public:project:" + slug)
		if err != nil {
			return repository.PublicVersion{}, err
		}
		versions = append(versions, v)
	}
	return repository.MergePublicVersions(versions...), nil
}

func (s *ProjectService) publicProjectVersion(key string) (repository.PublicVersion, error) {
	slug := publicSlug(key, "public:project:")
	query := `
SELECT MAX(updated_at) AS updated_at, COUNT(*) AS count FROM (
  SELECT projects.updated_at AS updated_at FROM projects WHERE projects.slug = ? AND projects.deleted_at IS NULL
  UNION ALL SELECT requirements.updated_at FROM requirements JOIN projects ON projects.id = requirements.project_id WHERE projects.slug = ? AND projects.deleted_at IS NULL
  UNION ALL SELECT cost_items.updated_at FROM cost_items JOIN projects ON projects.id = cost_items.project_id WHERE projects.slug = ? AND projects.deleted_at IS NULL
  UNION ALL SELECT timeline_phases.updated_at FROM timeline_phases JOIN projects ON projects.id = timeline_phases.project_id WHERE projects.slug = ? AND projects.deleted_at IS NULL
  UNION ALL SELECT milestones.updated_at FROM milestones JOIN projects ON projects.id = milestones.project_id WHERE projects.slug = ? AND projects.deleted_at IS NULL
  UNION ALL SELECT project_advantages.updated_at FROM project_advantages JOIN projects ON projects.id = project_advantages.project_id WHERE projects.slug = ? AND projects.deleted_at IS NULL
  UNION ALL SELECT compare_configs.updated_at FROM compare_configs JOIN projects ON projects.id = compare_configs.project_id WHERE projects.slug = ? AND projects.deleted_at IS NULL
  UNION ALL SELECT cases.updated_at FROM cases JOIN projects ON projects.id = cases.project_id WHERE projects.slug = ? AND projects.deleted_at IS NULL
  UNION ALL SELECT testimonials.updated_at FROM testimonials JOIN projects ON projects.id = testimonials.project_id WHERE projects.slug = ? AND projects.deleted_at IS NULL
  UNION ALL SELECT faqs.updated_at FROM faqs JOIN projects ON projects.id = faqs.project_id WHERE projects.slug = ? AND projects.deleted_at IS NULL
  UNION ALL SELECT project_news.created_at FROM project_news JOIN projects ON projects.id = project_news.project_id WHERE projects.slug = ? AND projects.deleted_at IS NULL
  UNION ALL SELECT pages.updated_at FROM pages JOIN project_news ON project_news.page_id = pages.id JOIN projects ON projects.id = project_news.project_id WHERE projects.slug = ? AND projects.deleted_at IS NULL AND pages.deleted_at IS NULL
) AS versions`
	return s.versionRepo.VersionFromQuery(query, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug)
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
	project.CoverImageVariants = ResolveImageVariants(project.CoverImage, UploadContextProject)
	return project, nil
}

// List returns projects with optional filtering and pagination.
func (s *ProjectService) List(req dto.ProjectListRequest) ([]model.Project, int64, error) {
	projects, total, err := s.repo.FindAll(repository.ProjectFilter{
		Name:    req.Name,
		Status:  req.Status,
		Page:    req.Page,
		PerPage: req.PerPage,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list projects: %w", err)
	}
	for i := range projects {
		projects[i].CoverImageVariants = ResolveImageVariants(projects[i].CoverImage, UploadContextProject)
	}
	return projects, total, nil
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

// fieldBuilder maps a compare field key to its label and value extraction function.
type fieldBuilder struct {
	Label   string
	Extract func(model.Project) string
	ItemsFn func([]model.Project) [][]string
}

// compareFieldBuilders maps each compare field key to its builder.
var compareFieldBuilders = map[string]fieldBuilder{
	"investment_amount":  {Label: "投资金额", Extract: func(p model.Project) string { return p.InvestmentAmount }},
	"processing_period":  {Label: "办理周期", Extract: func(p model.Project) string { return p.ProcessingPeriod }},
	"target_crowd":       {Label: "适合人群", Extract: func(p model.Project) string { return p.TargetCrowd }},
	"country":            {Label: "国家", Extract: func(p model.Project) string { return p.Country }},
	"costs_total":        {Label: "费用总计", Extract: func(p model.Project) string { return p.CostsTotal }},
	"requirements_count": {Label: "申请条件", Extract: func(p model.Project) string { return joinRequirements(p.Requirements) }, ItemsFn: pluckRequirements},
	"timeline_steps":     {Label: "流程步骤数", Extract: func(p model.Project) string { return fmt.Sprintf("%d 个阶段", len(p.TimelinePhases)) }},
	"overview_text":      {Label: "项目介绍", Extract: func(p model.Project) string { return p.OverviewText }},
	"tagline":            {Label: "标语", Extract: func(p model.Project) string { return p.Tagline }},
	"policy_title":       {Label: "政策标题", Extract: func(p model.Project) string { return p.PolicyTitle }},
}

// CompareRows returns formatted comparison rows for N projects.
// If fields is empty, all default fields from config.CompareFields are returned.
// Otherwise, only the requested field keys are included (unknown keys are skipped).
func (s *ProjectService) CompareRows(slugs []string, fields []string) (*dto.CompareResult, error) {
	projects, err := s.Compare(slugs)
	if err != nil {
		return nil, err
	}
	if len(projects) < 2 {
		return nil, errors.New("需要至少两个项目进行对比")
	}

	projInfo := make([]dto.CompareProject, len(projects))
	for i, p := range projects {
		projInfo[i] = dto.CompareProject{Title: p.Name, Slug: p.Slug}
	}

	// Determine which field keys to use
	selectedKeys := fields
	if len(selectedKeys) == 0 {
		selectedKeys = make([]string, len(config.CompareFields))
		for i, f := range config.CompareFields {
			selectedKeys[i] = f.Key
		}
	}

	// Build rows for selected keys, skipping unknown keys
	rows := make([]dto.CompareRow, 0, len(selectedKeys))
	for _, key := range selectedKeys {
		builder, ok := compareFieldBuilders[key]
		if !ok {
			continue
		}
		row := dto.CompareRow{
			Label:  builder.Label,
			Values: pluck(projects, builder.Extract),
		}
		if builder.ItemsFn != nil {
			row.Items = builder.ItemsFn(projects)
		}
		rows = append(rows, row)
	}

	return &dto.CompareResult{Projects: projInfo, Rows: rows}, nil
}

func pluck(projects []model.Project, fn func(model.Project) string) []string {
	values := make([]string, len(projects))
	for i, p := range projects {
		values[i] = fn(p)
	}
	return values
}

func joinRequirements(reqs []model.Requirement) string {
	labels := requirementLabels(reqs)
	if len(labels) == 0 {
		return ""
	}
	return strings.Join(labels, "；")
}

func requirementLabels(reqs []model.Requirement) []string {
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
	return labels
}

func pluckRequirements(projects []model.Project) [][]string {
	items := make([][]string, len(projects))
	for i, p := range projects {
		items[i] = requirementLabels(p.Requirements)
	}
	return items
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
func (s *ProjectService) Update(id uint64, req dto.UpdateProjectRequest) (*model.Project, error) {
	if id == 0 {
		return nil, errors.New("project id is required")
	}

	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// 检查 slug 唯一性（如果 slug 变更了）
	if req.Slug != "" && req.Slug != existing.Slug {
		other, err := s.repo.FindBySlug(req.Slug)
		if err == nil && other != nil && other.ID != id {
			return nil, fmt.Errorf("slug %s 已被使用", req.Slug)
		}
		if s.navRepo != nil {
			count, err := s.navRepo.CountByProjectID(id)
			if err != nil {
				return nil, fmt.Errorf("failed to check navigation references: %w", err)
			}
			if count > 0 {
				return nil, fmt.Errorf("%d 个导航项引用了此项目，请先解除引用", count)
			}
		}
	}

	existing.Slug = req.Slug
	existing.Name = req.Name
	existing.Country = req.Country
	existing.FlagEmoji = req.FlagEmoji
	existing.Tagline = req.Tagline
	existing.InvestmentAmount = req.InvestmentAmount
	existing.InvestmentValue = req.InvestmentValue
	existing.InvestmentCurrency = req.InvestmentCurrency
	existing.ProcessingPeriod = req.ProcessingPeriod
	existing.TargetCrowd = req.TargetCrowd
	existing.OverviewTitle = req.OverviewTitle
	existing.OverviewText = req.OverviewText
	existing.PolicyTitle = req.PolicyTitle
	existing.PolicyText = req.PolicyText
	existing.CostsTotal = req.CostsTotal
	existing.CostsNote = req.CostsNote
	existing.CtaText = req.CtaText
	existing.HeroTitle = req.HeroTitle
	existing.HeroDesc = req.HeroDesc
	existing.HeroGradient = req.HeroGradient
	existing.CoverImage = req.CoverImage
	existing.SortOrder = req.SortOrder
	existing.Status = req.Status

	if err := s.repo.Update(existing); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}
	return existing, nil
}

// Delete performs a soft delete on a project by ID.
// ListNews returns news pages linked to a project.
func (s *ProjectService) ListNews(projectID uint64) ([]model.Page, error) {
	return s.repo.FindNews(projectID)
}

// AddNews links news pages to a project.
func (s *ProjectService) AddNews(projectID uint64, pageIDs []uint64) error {
	return s.repo.AddNews(projectID, pageIDs)
}

// RemoveNews unlinks a news page from a project.
func (s *ProjectService) RemoveNews(projectID, pageID uint64) error {
	return s.repo.RemoveNews(projectID, pageID)
}

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

	caseIDs, testimonialIDs := s.preDeleteCleanup(id)

	err := repository.Tx(func(txRepo *repository.Repository) error {
		if err := cascadeDeleteResources(txRepo, id); err != nil {
			return err
		}
		return txRepo.Project.Delete(id)
	})
	if errors.Is(err, repository.ErrTxNotReady) {
		s.cascadeDeleteProjectResources(id)
		if err := s.repo.Delete(id); err != nil {
			return fmt.Errorf("failed to delete project: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	s.postDeleteHomeConfigCleanup(id, caseIDs, testimonialIDs)
	return nil
}

// cascadeDeleteResources deletes all child resources of a project using the
// given Repository. Any error triggers a rollback in the calling transaction.
func cascadeDeleteResources(repo *repository.Repository, id uint64) error {
	if err := repo.Requirement.DeleteByProjectID(id); err != nil {
		return err
	}
	if err := repo.CostItem.DeleteByProjectID(id); err != nil {
		return err
	}
	if err := repo.TimelinePhase.DeleteByProjectID(id); err != nil {
		return err
	}
	if err := repo.Milestone.DeleteByProjectID(id); err != nil {
		return err
	}
	if err := repo.ProjectAdvantage.DeleteByProjectID(id); err != nil {
		return err
	}
	if err := repo.CompareConfig.DeleteByProjectID(id); err != nil {
		return err
	}
	if err := repo.Project.DeleteNewsByProjectID(id); err != nil {
		return err
	}
	if err := repo.Case.DeleteByProjectID(id); err != nil {
		return err
	}
	if err := repo.Testimonial.DeleteByProjectID(id); err != nil {
		return err
	}
	if err := repo.FAQ.DeleteByProjectID(id); err != nil {
		return err
	}
	return nil
}

func (s *ProjectService) preDeleteCleanup(id uint64) (caseIDs, testimonialIDs []uint64) {
	if s.caseRepo != nil {
		if cases, _, err := s.caseRepo.FindAll(repository.CaseFilter{ProjectID: &id}); err == nil {
			for _, c := range cases {
				caseIDs = append(caseIDs, c.ID)
			}
		}
	}
	if s.testimonialRepo != nil {
		if testimonials, _, err := s.testimonialRepo.FindAll(repository.TestimonialFilter{ProjectID: &id}); err == nil {
			for _, t := range testimonials {
				testimonialIDs = append(testimonialIDs, t.ID)
			}
		}
	}
	return
}

func (s *ProjectService) cascadeDeleteProjectResources(id uint64) {
	if s.requirementRepo != nil {
		if err := s.requirementRepo.DeleteByProjectID(id); err != nil {
			logging.Logger.Warn("cascade delete requirements failed", "project_id", id, "error", err)
		}
	}
	if s.costItemRepo != nil {
		if err := s.costItemRepo.DeleteByProjectID(id); err != nil {
			logging.Logger.Warn("cascade delete cost_items failed", "project_id", id, "error", err)
		}
	}
	if s.timelinePhaseRepo != nil {
		if err := s.timelinePhaseRepo.DeleteByProjectID(id); err != nil {
			logging.Logger.Warn("cascade delete timeline_phases failed", "project_id", id, "error", err)
		}
	}
	if s.milestoneRepo != nil {
		if err := s.milestoneRepo.DeleteByProjectID(id); err != nil {
			logging.Logger.Warn("cascade delete milestones failed", "project_id", id, "error", err)
		}
	}
	if s.advantageRepo != nil {
		if err := s.advantageRepo.DeleteByProjectID(id); err != nil {
			logging.Logger.Warn("cascade delete project_advantages failed", "project_id", id, "error", err)
		}
	}
	if s.compareConfigRepo != nil {
		if err := s.compareConfigRepo.DeleteByProjectID(id); err != nil {
			logging.Logger.Warn("cascade delete compare_config failed", "project_id", id, "error", err)
		}
	}
	if err := s.repo.DeleteNewsByProjectID(id); err != nil {
		logging.Logger.Warn("cascade delete project_news failed", "project_id", id, "error", err)
	}
	if s.caseRepo != nil {
		if err := s.caseRepo.DeleteByProjectID(id); err != nil {
			logging.Logger.Warn("cascade delete cases failed", "project_id", id, "error", err)
		}
	}
	if s.testimonialRepo != nil {
		if err := s.testimonialRepo.DeleteByProjectID(id); err != nil {
			logging.Logger.Warn("cascade delete testimonials failed", "project_id", id, "error", err)
		}
	}
	if s.faqRepo != nil {
		if err := s.faqRepo.DeleteByProjectID(id); err != nil {
			logging.Logger.Warn("cascade delete faqs failed", "project_id", id, "error", err)
		}
	}
}

func (s *ProjectService) postDeleteHomeConfigCleanup(id uint64, caseIDs, testimonialIDs []uint64) {
	if s.homeConfigSvc == nil {
		return
	}
	if len(caseIDs) > 0 || len(testimonialIDs) > 0 {
		if err := s.homeConfigSvc.CleanupFeaturedRefs(caseIDs, testimonialIDs); err != nil {
			logging.Logger.Warn("home_config: failed to clean up featured refs after cascade delete",
				"project_id", id, "case_ids", caseIDs, "testimonial_ids", testimonialIDs, "error", err)
		}
	}
	if err := s.homeConfigSvc.RemoveFeaturedProjectID(id); err != nil {
		logging.Logger.Warn("home_config: failed to clean up featured project ref after delete",
			"project_id", id, "error", err)
	}
}
