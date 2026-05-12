package service

import (
	"errors"
	"fmt"

	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"

	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
)

var caseContentSanitizer = func() *bluemonday.Policy {
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

type CaseService struct {
	repo repository.CaseRepository
}

func NewCaseService(repo repository.CaseRepository) *CaseService {
	return &CaseService{repo: repo}
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
	cases, err := s.repo.FindAll("")
	if err != nil {
		return nil, fmt.Errorf("failed to list cases: %w", err)
	}
	return cases, nil
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
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	cases, err := s.repo.FindAll(search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list cases: %w", err)
	}

	total := int64(len(cases))
	start := (page - 1) * perPage
	if start >= len(cases) {
		return []model.Case{}, total, nil
	}
	end := start + perPage
	if end > len(cases) {
		end = len(cases)
	}
	return cases[start:end], total, nil
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
	c.Content = caseContentSanitizer.Sanitize(c.Content)
	if err := s.repo.Create(c); err != nil {
		return nil, fmt.Errorf("failed to create case: %w", err)
	}
	return c, nil
}

func (s *CaseService) Update(id uint64, c *model.Case) (*model.Case, error) {
	if c == nil {
		return nil, errors.New("case is nil")
	}
	if id == 0 {
		return nil, errors.New("case id is required")
	}
	c.ID = id
	c.Content = caseContentSanitizer.Sanitize(c.Content)
	if err := s.repo.Update(c); err != nil {
		return nil, fmt.Errorf("failed to update case: %w", err)
	}
	return c, nil
}

func (s *CaseService) Delete(id uint64) error {
	if id == 0 {
		return errors.New("case id is required")
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete case: %w", err)
	}
	return nil
}

func (s *CaseService) HardDelete(id uint64) error {
	if id == 0 {
		return errors.New("case id is required")
	}
	if err := s.repo.HardDelete(id); err != nil {
		return fmt.Errorf("failed to hard delete case: %w", err)
	}
	return nil
}
