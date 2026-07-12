package service

import (
	"errors"
	"fmt"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

// PageService handles business logic for content pages.
type PageService struct {
	repo        repository.PageRepository
	navRepo     repository.NavigationRepository
	versionRepo *repository.PublicVersionRepo
}

// NewPageService creates a new PageService with the given dependencies.
func NewPageService(repo repository.PageRepository, navRepo repository.NavigationRepository) *PageService {
	return &PageService{repo: repo, navRepo: navRepo}
}

func (s *PageService) RegisterPublicVersions(reg *PublicVersionRegistry) {
	reg.Register("public:pages:list", func(string) (repository.PublicVersion, error) {
		return tableVersion(s.versionRepo, "pages", "deleted_at IS NULL AND status = ?", "published")
	})
	reg.Register("public:page:", func(key string) (repository.PublicVersion, error) {
		return tableVersion(s.versionRepo, "pages", "deleted_at IS NULL AND status = ? AND slug = ?", "published", publicSlug(key, "public:page:"))
	})
}

// GetBySlug returns a published page by its slug.
func (s *PageService) GetBySlug(slug string) (*model.Page, error) {
	if slug == "" {
		return nil, errors.New("slug is required")
	}
	page, err := s.repo.FindBySlugPublished(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get page by slug: %w", err)
	}
	page.CoverImageVariants = ResolveImageVariants(page.CoverImage, UploadContextPageCover)
	return page, nil
}

// GetBySlugPreview returns a page by slug regardless of status.
func (s *PageService) GetBySlugPreview(slug string) (*model.Page, error) {
	if slug == "" {
		return nil, errors.New("slug is required")
	}
	page, err := s.repo.FindBySlug(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get page by slug: %w", err)
	}
	return page, nil
}

// List returns pages with optional filtering and pagination.
func (s *PageService) List(req dto.PageListRequest) ([]model.Page, int64, error) {
	pages, total, err := s.repo.FindAll(repository.PageFilter{
		PageType: req.PageType,
		Title:    req.Title,
		Status:   req.Status,
		Page:     req.Page,
		PerPage:  req.PerPage,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list pages: %w", err)
	}
	for i := range pages {
		pages[i].CoverImageVariants = ResolveImageVariants(pages[i].CoverImage, UploadContextPageCover)
	}
	return pages, total, nil
}

// Search returns pages matching a search query.
func (s *PageService) Search(q string) ([]model.Page, error) {
	if q == "" {
		return nil, errors.New("search query is required")
	}

	pages, err := s.repo.Search(q)
	if err != nil {
		return nil, fmt.Errorf("failed to search pages: %w", err)
	}
	return pages, nil
}

// Create creates a new page, sanitizing the content field against XSS.
func (s *PageService) Create(page *model.Page) (*model.Page, error) {
	if page == nil {
		return nil, errors.New("page is nil")
	}
	if page.Title == "" {
		return nil, errors.New("page title is required")
	}
	if page.Slug == "" {
		return nil, errors.New("page slug is required")
	}
	page.ID = 0
	page.Content = HTMLSanitizer.Sanitize(page.Content)
	if err := s.repo.Create(page); err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}
	return page, nil
}

// Update updates an existing page, sanitizing the content field against XSS.
func (s *PageService) Update(id uint64, req dto.UpdatePageRequest) (*model.Page, error) {
	if id == 0 {
		return nil, errors.New("page id is required")
	}

	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("page not found: %w", err)
	}

	// 检查 slug 唯一性（如果 slug 变更了）
	if req.Slug != "" && req.Slug != existing.Slug {
		other, err := s.repo.FindBySlug(req.Slug)
		if err == nil && other != nil && other.ID != id {
			return nil, fmt.Errorf("slug %s 已被使用", req.Slug)
		}
	}

	// 如果页面被导航引用，不允许修改 slug
	if req.Slug != "" && req.Slug != existing.Slug && s.navRepo != nil {
		count, err := s.navRepo.CountByPageID(id)
		if err != nil {
			return nil, fmt.Errorf("failed to check navigation references: %w", err)
		}
		if count > 0 {
			return nil, fmt.Errorf("%d 个导航项引用了此页面，请先解除引用", count)
		}
	}

	if req.Slug != "" {
		existing.Slug = req.Slug
	}
	existing.ProjectID = req.ProjectID
	existing.Title = req.Title
	existing.Content = HTMLSanitizer.Sanitize(req.Content)
	existing.CoverImage = req.CoverImage
	existing.MetaTitle = req.MetaTitle
	existing.MetaDescription = req.MetaDescription
	existing.Template = req.Template
	existing.PageType = req.PageType
	existing.Status = req.Status
	existing.SortOrder = req.SortOrder
	existing.IsPinned = req.IsPinned

	if err := s.repo.Update(existing); err != nil {
		return nil, fmt.Errorf("failed to update page: %w", err)
	}
	return existing, nil
}

// Delete removes a page by ID.
func (s *PageService) Delete(id uint64) error {
	if id == 0 {
		return errors.New("page id is required")
	}
	if s.navRepo != nil {
		count, err := s.navRepo.CountByPageID(id)
		if err != nil {
			return fmt.Errorf("failed to check navigation references: %w", err)
		}
		if count > 0 {
			return fmt.Errorf("%d 个导航项引用了此页面，请先解除引用", count)
		}
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete page: %w", err)
	}
	return nil
}
