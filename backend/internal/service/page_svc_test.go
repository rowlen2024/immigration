package service

import (
	"errors"
	"strings"
	"testing"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
	"time"
)

// mockPageRepo implements repository.PageRepository for testing.
type mockPageRepo struct {
	findByIDFn            func(id uint64) (*model.Page, error)
	findBySlugFn          func(slug string) (*model.Page, error)
	findAllFn             func(filter repository.PageFilter) ([]model.Page, int64, error)
	findOptionsFn         func(filter repository.PageFilter) ([]repository.PageOptionRow, int64, error)
	findBySlugPublishedFn func(slug string) (*model.Page, error)
	createFn              func(page *model.Page) error
	updateFn              func(page *model.Page) error
	deleteFn              func(id uint64) error
	searchFn              func(keyword string) ([]model.Page, error)
}

func (m *mockPageRepo) FindByID(id uint64) (*model.Page, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return nil, errors.New("not found")
}

func (m *mockPageRepo) FindBySlug(slug string) (*model.Page, error) {
	if m.findBySlugFn != nil {
		return m.findBySlugFn(slug)
	}
	return nil, errors.New("not found")
}

func (m *mockPageRepo) FindAll(filter repository.PageFilter) ([]model.Page, int64, error) {
	if m.findAllFn != nil {
		return m.findAllFn(filter)
	}
	return nil, 0, nil
}

func (m *mockPageRepo) FindOptions(filter repository.PageFilter) ([]repository.PageOptionRow, int64, error) {
	if m.findOptionsFn != nil {
		return m.findOptionsFn(filter)
	}
	return nil, 0, nil
}

func (m *mockPageRepo) FindBySlugPublished(slug string) (*model.Page, error) {
	if m.findBySlugPublishedFn != nil {
		return m.findBySlugPublishedFn(slug)
	}
	return nil, errors.New("not found")
}

func (m *mockPageRepo) Create(page *model.Page) error {
	if m.createFn != nil {
		return m.createFn(page)
	}
	return nil
}

func (m *mockPageRepo) Update(page *model.Page) error {
	if m.updateFn != nil {
		return m.updateFn(page)
	}
	return nil
}

func (m *mockPageRepo) Delete(id uint64) error {
	if m.deleteFn != nil {
		return m.deleteFn(id)
	}
	return nil
}

func (m *mockPageRepo) Search(keyword string) ([]model.Page, error) {
	if m.searchFn != nil {
		return m.searchFn(keyword)
	}
	return nil, nil
}

func (m *mockPageRepo) FindAllCoverImages() ([]string, error)            { return nil, nil }
func (m *mockPageRepo) FindAllContents() ([]string, error)               { return nil, nil }
func (m *mockPageRepo) Count() (int64, error)                            { return 0, nil }
func (m *mockPageRepo) CountByRange(start, end time.Time) (int64, error) { return 0, nil }

func TestPageService_OptionsClampsPerPage(t *testing.T) {
	repo := &mockPageRepo{
		findOptionsFn: func(filter repository.PageFilter) ([]repository.PageOptionRow, int64, error) {
			if filter.PageType != "news" || filter.Status != "published" || filter.Title != "visa" {
				t.Fatalf("unexpected filter: %#v", filter)
			}
			if filter.Page != 1 || filter.PerPage != 500 {
				t.Fatalf("expected pagination 1/500, got %d/%d", filter.Page, filter.PerPage)
			}
			return []repository.PageOptionRow{{ID: 3, Slug: "visa-news", Title: "Visa News"}}, 1, nil
		},
	}
	svc := NewPageService(repo, nil)

	items, total, err := svc.Options(dto.PageOptionRequest{
		OptionPaginationRequest: dto.OptionPaginationRequest{Page: 1, PerPage: 999},
		PageType:                "news",
		Status:                  "published",
		Title:                   "visa",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 1 || len(items) != 1 || items[0].Slug != "visa-news" {
		t.Fatalf("unexpected options result: total=%d items=%v", total, items)
	}
}

func TestPage_GetBySlug_Success(t *testing.T) {
	expected := &model.Page{ID: 1, Title: "About Us", Slug: "about"}
	repo := &mockPageRepo{
		findBySlugPublishedFn: func(slug string) (*model.Page, error) {
			if slug == "about" {
				return expected, nil
			}
			return nil, errors.New("not found")
		},
	}

	svc := NewPageService(repo, nil)

	page, err := svc.GetBySlug("about")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if page.ID != 1 || page.Title != "About Us" {
		t.Errorf("unexpected page: %+v", page)
	}
}

func TestPage_GetBySlug_EmptySlug(t *testing.T) {
	repo := &mockPageRepo{}
	svc := NewPageService(repo, nil)

	_, err := svc.GetBySlug("")
	if err == nil {
		t.Fatal("expected error for empty slug")
	}
}

func TestPage_GetBySlug_NotFound(t *testing.T) {
	repo := &mockPageRepo{
		findBySlugPublishedFn: func(slug string) (*model.Page, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewPageService(repo, nil)

	_, err := svc.GetBySlug("nonexistent")
	if err == nil {
		t.Fatal("expected error for not found")
	}
}

func TestPage_Search_Success(t *testing.T) {
	samplePages := []model.Page{
		{ID: 1, Title: "Immigration Guide", Slug: "guide"},
		{ID: 2, Title: "Immigration FAQ", Slug: "faq"},
	}

	repo := &mockPageRepo{
		searchFn: func(keyword string) ([]model.Page, error) {
			return samplePages, nil
		},
	}

	svc := NewPageService(repo, nil)

	pages, err := svc.Search("Immigration")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(pages) != 2 {
		t.Errorf("expected 2 pages, got %d", len(pages))
	}
}

func TestPage_Search_EmptyQuery(t *testing.T) {
	repo := &mockPageRepo{}
	svc := NewPageService(repo, nil)

	_, err := svc.Search("")
	if err == nil {
		t.Fatal("expected error for empty search query")
	}
}

func TestPage_Create_XSSSanitization(t *testing.T) {
	var savedContent string
	repo := &mockPageRepo{
		createFn: func(page *model.Page) error {
			savedContent = page.Content
			page.ID = 1
			return nil
		},
	}

	svc := NewPageService(repo, nil)

	xssContent := `<p>Hello</p><script>alert("xss")</script><b>World</b>`
	page, err := svc.Create(&model.Page{Title: "Test", Slug: "test", Content: xssContent})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if strings.Contains(savedContent, "<script>") {
		t.Errorf("expected script tag to be sanitized, got: %s", savedContent)
	}
	if !strings.Contains(savedContent, "<p>") {
		t.Errorf("expected <p> tag to be preserved, got: %s", savedContent)
	}
	if page.ID != 1 {
		t.Errorf("expected ID 1, got %d", page.ID)
	}
}

func TestPage_Create_NilPage(t *testing.T) {
	repo := &mockPageRepo{}
	svc := NewPageService(repo, nil)

	_, err := svc.Create(nil)
	if err == nil {
		t.Fatal("expected error for nil page")
	}
}

func TestPage_Create_MissingTitle(t *testing.T) {
	repo := &mockPageRepo{}
	svc := NewPageService(repo, nil)

	_, err := svc.Create(&model.Page{Slug: "test-slug"})
	if err == nil {
		t.Fatal("expected error for missing title")
	}
}

func TestPage_Create_MissingSlug(t *testing.T) {
	repo := &mockPageRepo{}
	svc := NewPageService(repo, nil)

	_, err := svc.Create(&model.Page{Title: "Test Title"})
	if err == nil {
		t.Fatal("expected error for missing slug")
	}
}

func TestPage_Update_Success(t *testing.T) {
	var savedContent string
	repo := &mockPageRepo{
		findByIDFn: func(id uint64) (*model.Page, error) {
			return &model.Page{ID: id, Slug: "existing", Title: "Old", Content: ""}, nil
		},
		updateFn: func(page *model.Page) error {
			savedContent = page.Content
			return nil
		},
	}

	svc := NewPageService(repo, nil)

	page, err := svc.Update(1, dto.UpdatePageRequest{Title: "Updated", Slug: "updated", Content: "<p>Safe</p><script>bad</script>"})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if page.ID != 1 {
		t.Errorf("expected ID 1, got %d", page.ID)
	}
	if strings.Contains(savedContent, "<script>") {
		t.Errorf("expected script tag to be sanitized in update, got: %s", savedContent)
	}
}

func TestPage_Update_NilPage(t *testing.T) {
	repo := &mockPageRepo{}
	svc := NewPageService(repo, nil)

	_, err := svc.Update(1, dto.UpdatePageRequest{})
	if err == nil {
		t.Fatal("expected error for nil page in update")
	}
}

func TestPage_Update_ZeroID(t *testing.T) {
	repo := &mockPageRepo{}
	svc := NewPageService(repo, nil)

	_, err := svc.Update(0, dto.UpdatePageRequest{Title: "T", Slug: "s", Content: "c"})
	if err == nil {
		t.Fatal("expected error for zero id")
	}
}

func TestPage_Update_RepoError(t *testing.T) {
	repo := &mockPageRepo{
		updateFn: func(page *model.Page) error {
			return errors.New("db error")
		},
	}

	svc := NewPageService(repo, nil)

	_, err := svc.Update(1, dto.UpdatePageRequest{Title: "T", Slug: "s", Content: "c"})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestPage_List_Success(t *testing.T) {
	samplePages := []model.Page{
		{ID: 1, Title: "Page A", Slug: "page-a"},
		{ID: 2, Title: "Page B", Slug: "page-b"},
		{ID: 3, Title: "Page C", Slug: "page-c"},
	}

	repo := &mockPageRepo{
		findAllFn: func(filter repository.PageFilter) ([]model.Page, int64, error) {
			return samplePages, 3, nil
		},
	}

	svc := NewPageService(repo, nil)

	pages, _, err := svc.List(dto.PageListRequest{Status: "published"})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(pages) != 3 {
		t.Errorf("expected 3 pages, got %d", len(pages))
	}
}

func TestPage_List_Empty(t *testing.T) {
	repo := &mockPageRepo{
		findAllFn: func(filter repository.PageFilter) ([]model.Page, int64, error) {
			return []model.Page{}, 0, nil
		},
	}

	svc := NewPageService(repo, nil)

	pages, _, err := svc.List(dto.PageListRequest{Status: "published"})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(pages) != 0 {
		t.Errorf("expected 0 pages, got %d", len(pages))
	}
}

func TestPage_List_Error(t *testing.T) {
	repo := &mockPageRepo{
		findAllFn: func(filter repository.PageFilter) ([]model.Page, int64, error) {
			return nil, 0, errors.New("db error")
		},
	}

	svc := NewPageService(repo, nil)

	_, _, err := svc.List(dto.PageListRequest{Status: "published"})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestPage_List_Paginated(t *testing.T) {
	samplePages := []model.Page{
		{ID: 1, Title: "Page A", Slug: "page-a"},
		{ID: 2, Title: "Page B", Slug: "page-b"},
	}

	repo := &mockPageRepo{
		findAllFn: func(filter repository.PageFilter) ([]model.Page, int64, error) {
			return samplePages, 4, nil
		},
	}

	svc := NewPageService(repo, nil)

	pages, total, err := svc.List(dto.PageListRequest{
		PaginationRequest: dto.PaginationRequest{Page: 1, PerPage: 2},
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if total != 4 {
		t.Errorf("expected total 4, got %d", total)
	}
	if len(pages) != 2 {
		t.Errorf("expected 2 pages, got %d", len(pages))
	}
}

func TestPage_List_Filtered(t *testing.T) {
	repo := &mockPageRepo{
		findAllFn: func(filter repository.PageFilter) ([]model.Page, int64, error) {
			if filter.PageType != "news" {
				t.Errorf("expected page_type 'news', got '%s'", filter.PageType)
			}
			if filter.Title != "test" {
				t.Errorf("expected title 'test', got '%s'", filter.Title)
			}
			return []model.Page{}, 0, nil
		},
	}

	svc := NewPageService(repo, nil)

	_, _, err := svc.List(dto.PageListRequest{
		PageType: "news",
		Title:    "test",
		Status:   "published",
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestPage_Delete_Success(t *testing.T) {
	deleted := false
	repo := &mockPageRepo{
		deleteFn: func(id uint64) error {
			deleted = true
			return nil
		},
	}

	svc := NewPageService(repo, nil)

	err := svc.Delete(5)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !deleted {
		t.Error("expected Delete to be called on repo")
	}
}

func TestPage_Delete_ZeroID(t *testing.T) {
	repo := &mockPageRepo{}
	svc := NewPageService(repo, nil)

	err := svc.Delete(0)
	if err == nil {
		t.Fatal("expected error for zero id")
	}
}

func TestPage_Delete_RepoError(t *testing.T) {
	repo := &mockPageRepo{
		deleteFn: func(id uint64) error {
			return errors.New("db error")
		},
	}

	svc := NewPageService(repo, nil)

	err := svc.Delete(1)
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestPage_Create_RepoError(t *testing.T) {
	repo := &mockPageRepo{
		createFn: func(page *model.Page) error {
			return errors.New("db error")
		},
	}

	svc := NewPageService(repo, nil)

	_, err := svc.Create(&model.Page{Title: "Test", Slug: "test", Content: "content"})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestPage_Search_Error(t *testing.T) {
	repo := &mockPageRepo{
		searchFn: func(keyword string) ([]model.Page, error) {
			return nil, errors.New("search error")
		},
	}

	svc := NewPageService(repo, nil)

	_, err := svc.Search("query")
	if err == nil {
		t.Fatal("expected error from repo")
	}
}
