package service

import (
	"errors"
	"fmt"

	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"

	"github.com/microcosm-cc/bluemonday"
)

// PageService handles business logic for content pages.
type PageService struct {
	repo    repository.PageRepository
	navRepo repository.NavigationRepository
}

// NewPageService creates a new PageService with the given repository.
func NewPageService(repo repository.PageRepository) *PageService {
	return &PageService{repo: repo}
}

var sanitizer = func() *bluemonday.Policy {
	p := bluemonday.NewPolicy()
	p.AllowElements("h1", "h2", "h3", "h4", "h5", "h6", "p", "br", "hr",
		"ul", "ol", "li", "blockquote", "pre", "code", "strong", "em", "u", "s",
		"a", "img", "table", "thead", "tbody", "tr", "th", "td",
		"div", "span", "iframe", "video", "source")
	p.AllowAttrs("src", "alt", "title", "width", "height").OnElements("img")
	p.AllowAttrs("href", "title", "target", "rel").OnElements("a")
	p.AllowAttrs("src", "frameborder", "allowfullscreen").OnElements("iframe")
	p.AllowAttrs("src", "controls", "width", "height").OnElements("video", "source")
	p.AllowAttrs("style", "class").OnElements("span", "div", "td", "th")
	p.AllowAttrs("class").OnElements("table", "thead", "tbody", "tr", "img", "a")
	p.AllowStyles("color", "background-color", "text-align").OnElements("span", "td", "th")
	p.AllowURLSchemes("http", "https", "mailto")
	p.AllowRelativeURLs(true)
	p.RequireNoFollowOnLinks(true)
	return p
}()

// GetBySlug returns a published page by its slug.
func (s *PageService) GetBySlug(slug string) (*model.Page, error) {
	if slug == "" {
		return nil, errors.New("slug is required")
	}
	page, err := s.repo.FindBySlugPublished(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get page by slug: %w", err)
	}
	return page, nil
}

// List returns all published pages.
func (s *PageService) List() ([]model.Page, error) {
	pages, err := s.repo.FindAllPublished()
	if err != nil {
		return nil, fmt.Errorf("failed to list pages: %w", err)
	}
	return pages, nil
}

// AdminList returns paginated pages, optionally filtered by page_type, search, and status.
func (s *PageService) AdminList(page, perPage int, pageType, search, status string) ([]model.Page, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	pages, err := s.repo.FindAll(pageType, search, status)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list pages: %w", err)
	}

	total := int64(len(pages))
	start := (page - 1) * perPage
	if start >= len(pages) {
		return []model.Page{}, total, nil
	}
	end := start + perPage
	if end > len(pages) {
		end = len(pages)
	}
	return pages[start:end], total, nil
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
	page.Content = sanitizer.Sanitize(page.Content)
	if err := s.repo.Create(page); err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}
	return page, nil
}

// Update updates an existing page, sanitizing the content field against XSS.
func (s *PageService) Update(id uint64, page *model.Page) (*model.Page, error) {
	if page == nil {
		return nil, errors.New("page is nil")
	}
	if id == 0 {
		return nil, errors.New("page id is required")
	}

	if s.navRepo != nil {
		count, err := s.navRepo.CountByPageID(id)
		if err != nil {
			return nil, fmt.Errorf("failed to check navigation references: %w", err)
		}
		if count > 0 {
			oldBySlug, slugErr := s.repo.FindBySlug(page.Slug)
			if slugErr != nil || oldBySlug == nil || oldBySlug.ID != id {
				return nil, fmt.Errorf("%d 个导航项引用了此页面，请先解除引用", count)
			}
		}
	}

	page.ID = id
	page.Content = sanitizer.Sanitize(page.Content)
	if err := s.repo.Update(page); err != nil {
		return nil, fmt.Errorf("failed to update page: %w", err)
	}
	return page, nil
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
