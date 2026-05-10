package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// handlerMockPageRepo implements repository.PageRepository.
type handlerMockPageRepo struct {
	findBySlug func(slug string) (*model.Page, error)
}

func (m *handlerMockPageRepo) FindBySlug(slug string) (*model.Page, error) {
	return m.findBySlug(slug)
}
func (m *handlerMockPageRepo) FindAll(pageType, search, status string) ([]model.Page, error) { return nil, nil }
func (m *handlerMockPageRepo) FindAllPublished() ([]model.Page, error)           { return nil, nil }
func (m *handlerMockPageRepo) FindBySlugPublished(slug string) (*model.Page, error) {
	return m.findBySlug(slug)
}
func (m *handlerMockPageRepo) FindByProjectID(projectID uint64) ([]model.Page, error) {
	return nil, nil
}
func (m *handlerMockPageRepo) Create(page *model.Page) error { return nil }
func (m *handlerMockPageRepo) Update(page *model.Page) error { return nil }
func (m *handlerMockPageRepo) Delete(id uint64) error        { return nil }
func (m *handlerMockPageRepo) Search(keyword string) ([]model.Page, error) { return nil, nil }

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

	pageSvc := service.NewPageService(mockRepo)
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

	pageSvc := service.NewPageService(mockRepo)
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
