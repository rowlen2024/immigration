package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrProjectCategoryName       = errors.New("请输入类目名称")
	ErrProjectCategorySlug       = errors.New("类目 slug 只能由小写字母、数字和单个连字符组成")
	ErrProjectCategorySlugTaken  = errors.New("类目 slug 已被使用")
	ErrProjectCategorySlugFrozen = errors.New("类目首次发布后不能修改 slug")
	ErrProjectCategoryInUse      = errors.New("已有导航项引用此类目，请先解除引用")
	projectCategorySlug          = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
)

type ProjectCategoryService struct {
	repo        *repository.ProjectCategoryRepo
	versionRepo *repository.PublicVersionRepo
}

func (s *ProjectCategoryService) List(req dto.ProjectCategoryListRequest) ([]repository.ProjectCategoryAdminRow, int64, error) {
	return s.repo.FindAll(repository.ProjectCategoryFilter{
		Name: req.Name, Status: req.Status, Page: req.Page, PerPage: req.PerPage,
	})
}

func (s *ProjectCategoryService) Options(req dto.ProjectCategoryListRequest) ([]dto.ProjectCategoryOption, int64, error) {
	categories, total, err := s.repo.FindOptions(repository.ProjectCategoryFilter{
		Name: req.Name, Status: req.Status, Page: req.Page, PerPage: req.PerPage,
	})
	if err != nil {
		return nil, 0, err
	}
	options := make([]dto.ProjectCategoryOption, 0, len(categories))
	for _, category := range categories {
		options = append(options, dto.ProjectCategoryOption{
			ID: category.ID, Name: category.Name, Slug: category.Slug, Status: category.Status,
		})
	}
	return options, total, nil
}

func (s *ProjectCategoryService) Create(req dto.SaveProjectCategoryRequest) (*model.ProjectCategory, error) {
	category := &model.ProjectCategory{Status: model.ProjectCategoryDraft}
	if err := applyProjectCategoryRequest(category, req); err != nil {
		return nil, err
	}
	if err := s.repo.Create(category); err != nil {
		return nil, categoryWriteError(err)
	}
	return category, nil
}

func (s *ProjectCategoryService) Update(id uint64, req dto.SaveProjectCategoryRequest) (*model.ProjectCategory, error) {
	category, err := s.repo.Update(id, func(existing *model.ProjectCategory) error {
		return applyProjectCategoryRequest(existing, req)
	})
	return category, categoryWriteError(err)
}

func applyProjectCategoryRequest(category *model.ProjectCategory, req dto.SaveProjectCategoryRequest) error {
	name := strings.TrimSpace(req.Name)
	slug := strings.TrimSpace(req.Slug)
	if name == "" {
		return ErrProjectCategoryName
	}
	if !projectCategorySlug.MatchString(slug) {
		return ErrProjectCategorySlug
	}
	wasPublished := category.PublishedAt != nil || category.Status == model.ProjectCategoryPublished
	if wasPublished && category.Slug != slug {
		return ErrProjectCategorySlugFrozen
	}
	category.Name = name
	category.Slug = slug
	category.Intro = strings.TrimSpace(req.Intro)
	category.SEOTitle = strings.TrimSpace(req.SEOTitle)
	category.SEODescription = strings.TrimSpace(req.SEODescription)
	if req.Status != "" {
		category.Status = req.Status
	}
	if category.PublishedAt == nil && (wasPublished || category.Status == model.ProjectCategoryPublished) {
		now := time.Now()
		category.PublishedAt = &now
	}
	return nil
}

func categoryWriteError(err error) error {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return ErrProjectCategorySlugTaken
	}
	return err
}

func (s *ProjectCategoryService) Delete(id uint64) error {
	err := s.repo.Delete(id)
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1451 {
		return ErrProjectCategoryInUse
	}
	return err
}

func (s *ProjectCategoryService) Projects(id uint64, page, perPage int) ([]repository.ProjectCategoryProjectRow, int64, error) {
	if _, err := s.repo.FindByID(id); err != nil {
		return nil, 0, err
	}
	return s.repo.FindProjects(id, page, perPage)
}

func (s *ProjectCategoryService) UpdateProjectOrder(id, projectID uint64, sortOrder int32) error {
	return s.repo.UpdateProjectOrder(id, projectID, sortOrder)
}

func (s *ProjectCategoryService) PublicList() ([]dto.ProjectCategoryPublicSummary, error) {
	rows, err := s.repo.FindPublicSummaries()
	if err != nil {
		return nil, err
	}
	summaries := make([]dto.ProjectCategoryPublicSummary, 0, len(rows))
	for _, row := range rows {
		summaries = append(summaries, dto.ProjectCategoryPublicSummary(row))
	}
	return summaries, nil
}

func (s *ProjectCategoryService) PublicDetail(slug string) (*dto.ProjectCategoryPublicDetail, error) {
	category, err := s.repo.FindPublishedBySlug(slug)
	if err != nil {
		return nil, err
	}
	projects, err := s.repo.FindPublishedProjects(category.ID)
	if err != nil {
		return nil, err
	}
	detail := &dto.ProjectCategoryPublicDetail{
		ProjectCategoryPublicSummary: dto.ProjectCategoryPublicSummary{
			ID: category.ID, Name: category.Name, Slug: category.Slug,
			UpdatedAt: category.UpdatedAt, ProjectCount: int64(len(projects)),
		},
		Intro: category.Intro, SEOTitle: category.SEOTitle, SEODescription: category.SEODescription,
		Projects: make([]dto.ProjectCategoryCard, 0, len(projects)),
	}
	for _, project := range projects {
		if project.UpdatedAt.After(detail.UpdatedAt) {
			detail.UpdatedAt = project.UpdatedAt
		}
		detail.Projects = append(detail.Projects, dto.ProjectCategoryCard{
			ID: project.ID, Name: project.Name, Slug: project.Slug, Country: project.Country, FlagEmoji: project.FlagEmoji,
			Tagline: project.Tagline, OverviewText: project.OverviewText,
			InvestmentAmount: project.InvestmentAmount, ProcessingPeriod: project.ProcessingPeriod,
			TargetCrowd: project.TargetCrowd, CoverImage: project.CoverImage,
			CoverImageVariants: ResolveImageVariants(project.CoverImage, UploadContextProject),
		})
	}
	return detail, nil
}

func (s *ProjectCategoryService) RegisterPublicVersions(reg *PublicVersionRegistry) {
	reg.Register("public:project-categories:list", s.publicVersion)
	reg.Register("public:project-categories:detail:", s.publicVersion)
}

func (s *ProjectCategoryService) publicVersion(key string) (repository.PublicVersion, error) {
	const sql = `SELECT MAX(updated_at) AS updated_at, COUNT(*) AS count FROM (
		SELECT c.updated_at FROM project_categories c WHERE c.status = 'published' %s
		UNION ALL
		SELECT p.updated_at FROM project_categories c
		JOIN project_category_memberships m ON m.category_id = c.id
		JOIN projects p ON p.id = m.project_id AND p.status = 1 AND p.deleted_at IS NULL
		WHERE c.status = 'published' %s
	) AS versions`
	const detailPrefix = "public:project-categories:detail:"
	if strings.HasPrefix(key, detailPrefix) {
		slug := publicSlug(key, detailPrefix)
		return s.versionRepo.VersionFromQuery(fmt.Sprintf(sql, "AND c.slug = ?", "AND c.slug = ?"), slug, slug)
	}
	return s.versionRepo.VersionFromQuery(fmt.Sprintf(sql, "", ""))
}
