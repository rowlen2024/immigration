package service

import (
	"errors"
	"testing"

	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
	"time"
)

// searchMockFAQRepo implements repository.FAQRepository for search service tests.
type searchMockFAQRepo struct {
	searchFn func(keyword string) ([]model.FAQ, error)
}

func (m *searchMockFAQRepo) FindByID(id uint64) (*model.FAQ, error)         { return nil, nil }
func (m *searchMockFAQRepo) FindDistinctProjects() ([]model.Project, error) { return nil, nil }
func (m *searchMockFAQRepo) FindAll(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
	return nil, 0, nil
}
func (m *searchMockFAQRepo) Create(faq *model.FAQ) error              { return nil }
func (m *searchMockFAQRepo) Update(faq *model.FAQ) error              { return nil }
func (m *searchMockFAQRepo) Delete(id uint64) error                   { return nil }
func (m *searchMockFAQRepo) DeleteByProjectID(projectID uint64) error { return nil }
func (m *searchMockFAQRepo) Search(keyword string) ([]model.FAQ, error) {
	if m.searchFn != nil {
		return m.searchFn(keyword)
	}
	return nil, nil
}

// searchMockPageRepo implements repository.PageRepository for search service tests.
type searchMockPageRepo struct {
	searchFn func(keyword string) ([]model.Page, error)
}

func (m *searchMockPageRepo) FindByID(id uint64) (*model.Page, error)     { return nil, nil }
func (m *searchMockPageRepo) FindBySlug(slug string) (*model.Page, error) { return nil, nil }
func (m *searchMockPageRepo) FindAll(filter repository.PageFilter) ([]model.Page, int64, error) {
	return nil, 0, nil
}
func (m *searchMockPageRepo) FindOptions(filter repository.PageFilter) ([]repository.PageOptionRow, int64, error) {
	return nil, 0, nil
}
func (m *searchMockPageRepo) FindBySlugPublished(slug string) (*model.Page, error) {
	return nil, nil
}
func (m *searchMockPageRepo) FindRelatedBySlug(string, int) ([]model.Page, error) { return nil, nil }
func (m *searchMockPageRepo) FindProjectsByPageID(uint64) ([]model.PageProject, error) {
	return nil, nil
}
func (m *searchMockPageRepo) Create(page *model.Page) error { return nil }
func (m *searchMockPageRepo) Update(page *model.Page) error { return nil }
func (m *searchMockPageRepo) Delete(id uint64) error        { return nil }
func (m *searchMockPageRepo) Search(keyword string) ([]model.Page, error) {
	if m.searchFn != nil {
		return m.searchFn(keyword)
	}
	return nil, nil
}

func TestSearch_Success(t *testing.T) {
	faqRepo := &searchMockFAQRepo{
		searchFn: func(keyword string) ([]model.FAQ, error) {
			return []model.FAQ{
				{ID: 1, Question: "How to apply?", Answer: "Process details here"},
			}, nil
		},
	}
	pageRepo := &searchMockPageRepo{
		searchFn: func(keyword string) ([]model.Page, error) {
			return []model.Page{
				{ID: 1, Title: "Application Guide", Slug: "application-guide"},
				{ID: 2, Title: "FAQ Page", Slug: "faq"},
			}, nil
		},
	}

	svc := &SearchService{faqRepo: faqRepo, pageRepo: pageRepo}

	results, err := svc.Search("apply")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(results.FAQs) != 1 {
		t.Errorf("expected 1 FAQ result, got %d", len(results.FAQs))
	}
	if len(results.Pages) != 2 {
		t.Errorf("expected 2 page results, got %d", len(results.Pages))
	}
}

func TestSearch_EmptyResults(t *testing.T) {
	faqRepo := &searchMockFAQRepo{
		searchFn: func(keyword string) ([]model.FAQ, error) {
			return []model.FAQ{}, nil
		},
	}
	pageRepo := &searchMockPageRepo{
		searchFn: func(keyword string) ([]model.Page, error) {
			return []model.Page{}, nil
		},
	}

	svc := &SearchService{faqRepo: faqRepo, pageRepo: pageRepo}

	results, err := svc.Search("nonexistent")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(results.FAQs) != 0 {
		t.Errorf("expected 0 FAQ results, got %d", len(results.FAQs))
	}
	if len(results.Pages) != 0 {
		t.Errorf("expected 0 page results, got %d", len(results.Pages))
	}
}

func TestSearch_FAQRepoError(t *testing.T) {
	faqRepo := &searchMockFAQRepo{
		searchFn: func(keyword string) ([]model.FAQ, error) {
			return nil, errors.New("faq search error")
		},
	}
	pageRepo := &searchMockPageRepo{
		searchFn: func(keyword string) ([]model.Page, error) {
			return nil, nil
		},
	}

	svc := &SearchService{faqRepo: faqRepo, pageRepo: pageRepo}

	_, err := svc.Search("test")
	if err == nil {
		t.Fatal("expected error from faq repo")
	}
}

func TestSearch_PageRepoError(t *testing.T) {
	faqRepo := &searchMockFAQRepo{
		searchFn: func(keyword string) ([]model.FAQ, error) {
			return []model.FAQ{}, nil
		},
	}
	pageRepo := &searchMockPageRepo{
		searchFn: func(keyword string) ([]model.Page, error) {
			return nil, errors.New("page search error")
		},
	}

	svc := &SearchService{faqRepo: faqRepo, pageRepo: pageRepo}

	_, err := svc.Search("test")
	if err == nil {
		t.Fatal("expected error from page repo")
	}
}

func TestSearch_SpecialCharacters(t *testing.T) {
	faqRepo := &searchMockFAQRepo{
		searchFn: func(keyword string) ([]model.FAQ, error) {
			if keyword != "test & query <script>" {
				t.Errorf("expected keyword 'test & query <script>', got '%s'", keyword)
			}
			return []model.FAQ{}, nil
		},
	}
	pageRepo := &searchMockPageRepo{
		searchFn: func(keyword string) ([]model.Page, error) {
			return []model.Page{}, nil
		},
	}

	svc := &SearchService{faqRepo: faqRepo, pageRepo: pageRepo}

	_, err := svc.Search("test & query <script>")
	if err != nil {
		t.Fatalf("expected success with special chars, got error: %v", err)
	}
}

func TestSearch_KeywordForwardedCorrectly(t *testing.T) {
	var faqKeyword, pageKeyword string
	faqRepo := &searchMockFAQRepo{
		searchFn: func(keyword string) ([]model.FAQ, error) {
			faqKeyword = keyword
			return []model.FAQ{}, nil
		},
	}
	pageRepo := &searchMockPageRepo{
		searchFn: func(keyword string) ([]model.Page, error) {
			pageKeyword = keyword
			return []model.Page{}, nil
		},
	}

	svc := &SearchService{faqRepo: faqRepo, pageRepo: pageRepo}

	_, err := svc.Search("immigration")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if faqKeyword != "immigration" {
		t.Errorf("FAQ repo received '%s', expected 'immigration'", faqKeyword)
	}
	if pageKeyword != "immigration" {
		t.Errorf("Page repo received '%s', expected 'immigration'", pageKeyword)
	}
}

func (m *searchMockPageRepo) FindAllCoverImages() ([]string, error)            { return nil, nil }
func (m *searchMockPageRepo) FindAllContents() ([]string, error)               { return nil, nil }
func (m *searchMockPageRepo) Count() (int64, error)                            { return 0, nil }
func (m *searchMockPageRepo) CountByRange(start, end time.Time) (int64, error) { return 0, nil }
