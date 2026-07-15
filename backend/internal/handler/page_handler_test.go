package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
	"mygo-immigration/backend/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// handlerMockPageRepo implements repository.PageRepository.
type handlerMockPageRepo struct {
	findByIDFn          func(id uint64) (*model.Page, error)
	findBySlug          func(slug string) (*model.Page, error)
	createFn            func(page *model.Page) error
	updateFn            func(page *model.Page) error
	findRelatedBySlugFn func(slug string, limit int) ([]model.Page, error)
}

func (m *handlerMockPageRepo) FindByID(id uint64) (*model.Page, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return nil, nil
}

func (m *handlerMockPageRepo) FindBySlug(slug string) (*model.Page, error) {
	return m.findBySlug(slug)
}
func (m *handlerMockPageRepo) FindAll(filter repository.PageFilter) ([]model.Page, int64, error) {
	return nil, 0, nil
}
func (m *handlerMockPageRepo) FindOptions(filter repository.PageFilter) ([]repository.PageOptionRow, int64, error) {
	return nil, 0, nil
}
func (m *handlerMockPageRepo) FindBySlugPublished(slug string) (*model.Page, error) {
	return m.findBySlug(slug)
}
func (m *handlerMockPageRepo) FindRelatedBySlug(slug string, limit int) ([]model.Page, error) {
	if m.findRelatedBySlugFn != nil {
		return m.findRelatedBySlugFn(slug, limit)
	}
	return nil, nil
}
func (m *handlerMockPageRepo) FindProjectsByPageID(uint64) ([]model.PageProject, error) {
	return nil, nil
}
func (m *handlerMockPageRepo) Create(page *model.Page) error {
	if m.createFn != nil {
		return m.createFn(page)
	}
	return nil
}
func (m *handlerMockPageRepo) Update(page *model.Page) error {
	if m.updateFn != nil {
		return m.updateFn(page)
	}
	return nil
}
func (m *handlerMockPageRepo) Delete(id uint64) error                           { return nil }
func (m *handlerMockPageRepo) Search(keyword string) ([]model.Page, error)      { return nil, nil }
func (m *handlerMockPageRepo) FindAllCoverImages() ([]string, error)            { return nil, nil }
func (m *handlerMockPageRepo) FindAllContents() ([]string, error)               { return nil, nil }
func (m *handlerMockPageRepo) Count() (int64, error)                            { return 0, nil }
func (m *handlerMockPageRepo) CountByRange(start, end time.Time) (int64, error) { return 0, nil }

func TestPageHandler_GetPage_Success(t *testing.T) {
	mockRepo := &handlerMockPageRepo{
		findBySlug: func(slug string) (*model.Page, error) {
			return &model.Page{
				ID:      1,
				Title:   "About Us",
				Slug:    slug,
				Content: "<p>About content</p>",
				Status:  "published",
			}, nil
		},
	}

	pageSvc := service.NewPageService(mockRepo, nil)
	h := &Handler{svc: &service.Service{Page: pageSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/pages/about")
	c.Params = gin.Params{{Key: "slug", Value: "about"}}

	h.GetPage(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var resp dto.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Code != 200 {
		t.Errorf("expected response code 200, got %d", resp.Code)
	}
}

func TestPageHandler_GetPage_NotFound(t *testing.T) {
	mockRepo := &handlerMockPageRepo{
		findBySlug: func(slug string) (*model.Page, error) {
			return nil, errors.New("record not found")
		},
	}

	pageSvc := service.NewPageService(mockRepo, nil)
	h := &Handler{svc: &service.Service{Page: pageSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/pages/nonexistent")
	c.Params = gin.Params{{Key: "slug", Value: "nonexistent"}}

	h.GetPage(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestPageHandler_GetPage_MissingSlug(t *testing.T) {
	h := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/pages/")

	h.GetPage(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestPageHandler_GetRelatedPages(t *testing.T) {
	mockRepo := &handlerMockPageRepo{findRelatedBySlugFn: func(slug string, limit int) ([]model.Page, error) {
		if slug != "current-news" || limit != 4 {
			t.Fatalf("unexpected query: slug=%q limit=%d", slug, limit)
		}
		return []model.Page{{ID: 2, Title: "Related", Slug: "related"}}, nil
	}}
	h := &Handler{svc: &service.Service{Page: service.NewPageService(mockRepo, nil)}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/related-pages?slug=current-news")
	h.GetRelatedPages(c)

	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"slug":"related"`) {
		t.Fatalf("unexpected response: status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestPageHandler_GetRelatedPages_MissingSlug(t *testing.T) {
	h := &Handler{}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/related-pages")

	h.GetRelatedPages(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestPageHandler_GetRelatedPages_EmptyArray(t *testing.T) {
	mockRepo := &handlerMockPageRepo{findRelatedBySlugFn: func(string, int) ([]model.Page, error) {
		return []model.Page{}, nil
	}}
	h := &Handler{svc: &service.Service{Page: service.NewPageService(mockRepo, nil)}}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/related-pages?slug=current-news")

	h.GetRelatedPages(c)

	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"data":[]`) {
		t.Fatalf("expected HTTP 200 with empty data array, got status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestPageHandler_GetRelatedPages_NotFound(t *testing.T) {
	mockRepo := &handlerMockPageRepo{findRelatedBySlugFn: func(string, int) ([]model.Page, error) {
		return nil, gorm.ErrRecordNotFound
	}}
	h := &Handler{svc: &service.Service{Page: service.NewPageService(mockRepo, nil)}}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/related-pages?slug=missing")

	h.GetRelatedPages(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestPageHandler_GetRelatedPages_RepositoryError(t *testing.T) {
	mockRepo := &handlerMockPageRepo{findRelatedBySlugFn: func(string, int) ([]model.Page, error) {
		return nil, errors.New("query failed")
	}}
	h := &Handler{svc: &service.Service{Page: service.NewPageService(mockRepo, nil)}}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/related-pages?slug=current-news")

	h.GetRelatedPages(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestPageHandler_CreatePage_BindsPinned(t *testing.T) {
	var createdPinned bool
	mockRepo := &handlerMockPageRepo{
		findBySlug: func(string) (*model.Page, error) { return nil, errors.New("not found") },
		createFn: func(page *model.Page) error {
			createdPinned = page.IsPinned
			page.ID = 99
			return nil
		},
	}

	pageSvc := service.NewPageService(mockRepo, nil)
	h := &Handler{svc: &service.Service{Page: pageSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makePostRequest("/api/v1/admin/pages", dto.CreatePageRequest{
		Title:    "New Page",
		Slug:     "new-page",
		IsPinned: true,
	})

	h.CreatePage(c)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d, body: %s", w.Code, w.Body.String())
	}
	if !createdPinned {
		t.Error("expected create handler to bind is_pinned")
	}
	if !strings.Contains(w.Body.String(), `"is_pinned":true`) {
		t.Errorf("expected response to include is_pinned, body: %s", w.Body.String())
	}
}

func TestPageHandler_CreatePage_IgnoresServerManagedFields(t *testing.T) {
	var received *model.Page
	mockRepo := &handlerMockPageRepo{
		createFn: func(page *model.Page) error {
			copyPage := *page
			received = &copyPage
			page.ID = 99
			return nil
		},
	}

	pageSvc := service.NewPageService(mockRepo, nil)
	h := &Handler{svc: &service.Service{Page: pageSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makePostRequest("/api/v1/admin/pages", map[string]interface{}{
		"id":               123,
		"title":            "adadfsfsfsdgdhfjh",
		"slug":             "adadfsfsfsdgdhfjh",
		"content":          "<p>ggggggggggggggggggggggggdfgdfytjttttttttttttttttttttttttttttttttttttt</p>",
		"cover_image":      "/uploads/1784084792587275333.png",
		"meta_title":       "",
		"meta_description": "",
		"template":         "default",
		"page_type":        "default",
		"status":           "draft",
		"is_pinned":        false,
		"tags":             []string{},
		"created_at":       "",
		"updated_at":       "",
	})

	h.CreatePage(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d, body: %s", w.Code, w.Body.String())
	}
	if received == nil {
		t.Fatal("expected repository to receive page")
	}
	if received.ID != 0 || !received.CreatedAt.IsZero() || !received.UpdatedAt.IsZero() {
		t.Fatalf("expected server-managed fields to remain zero, got ID=%d created_at=%v updated_at=%v", received.ID, received.CreatedAt, received.UpdatedAt)
	}
}

func TestPageHandler_UpdatePage_BindsPinned(t *testing.T) {
	var updatedPinned bool
	mockRepo := &handlerMockPageRepo{
		findByIDFn: func(id uint64) (*model.Page, error) {
			return &model.Page{ID: id, Title: "Old", Slug: "old"}, nil
		},
		findBySlug: func(string) (*model.Page, error) { return nil, errors.New("not found") },
		updateFn: func(page *model.Page) error {
			updatedPinned = page.IsPinned
			return nil
		},
	}

	pageSvc := service.NewPageService(mockRepo, nil)
	h := &Handler{svc: &service.Service{Page: pageSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makePostRequest("/api/v1/admin/pages/5", dto.UpdatePageRequest{
		Title:    "Updated Page",
		Slug:     "updated-page",
		IsPinned: true,
	})
	c.Params = gin.Params{{Key: "id", Value: "5"}}

	h.UpdatePage(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
	if !updatedPinned {
		t.Error("expected update handler to bind is_pinned")
	}
	if !strings.Contains(w.Body.String(), `"is_pinned":true`) {
		t.Errorf("expected response to include is_pinned, body: %s", w.Body.String())
	}
}
