package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
	"mygo-immigration/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// handlerMockFAQRepo implements repository.FAQRepository.
type handlerMockFAQRepo struct {
	findByIDFn             func(id uint64) (*model.FAQ, error)
	findAllFn              func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error)
	findDistinctProjectsFn func() ([]model.Project, error)
	createFn               func(faq *model.FAQ) error
	updateFn               func(faq *model.FAQ) error
	deleteFn               func(id uint64) error
	deleteByProjectIDFn    func(projectID uint64) error
	searchFn               func(keyword string) ([]model.FAQ, error)
}

func (m *handlerMockFAQRepo) FindByID(id uint64) (*model.FAQ, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return nil, nil
}
func (m *handlerMockFAQRepo) FindAll(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
	if m.findAllFn != nil {
		return m.findAllFn(params)
	}
	return nil, 0, nil
}
func (m *handlerMockFAQRepo) FindDistinctProjects() ([]model.Project, error) {
	if m.findDistinctProjectsFn != nil {
		return m.findDistinctProjectsFn()
	}
	return nil, nil
}
func (m *handlerMockFAQRepo) Create(faq *model.FAQ) error {
	if m.createFn != nil {
		return m.createFn(faq)
	}
	return nil
}
func (m *handlerMockFAQRepo) Update(faq *model.FAQ) error {
	if m.updateFn != nil {
		return m.updateFn(faq)
	}
	return nil
}
func (m *handlerMockFAQRepo) Delete(id uint64) error {
	if m.deleteFn != nil {
		return m.deleteFn(id)
	}
	return nil
}
func (m *handlerMockFAQRepo) DeleteByProjectID(projectID uint64) error {
	if m.deleteByProjectIDFn != nil {
		return m.deleteByProjectIDFn(projectID)
	}
	return nil
}
func (m *handlerMockFAQRepo) Search(keyword string) ([]model.FAQ, error) {
	if m.searchFn != nil {
		return m.searchFn(keyword)
	}
	return nil, nil
}

func TestFAQHandler_ListFAQs(t *testing.T) {
	mockRepo := &handlerMockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return []repository.FAQWithProject{
				{FAQ: model.FAQ{ID: 1, Question: "How to apply?", Answer: "Fill the form.", CreatedAt: time.Now(), UpdatedAt: time.Now()}},
				{FAQ: model.FAQ{ID: 2, Question: "Processing time?", Answer: "2-3 months.", CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			}, 2, nil
		},
	}

	faqSvc := service.NewFAQService(mockRepo)
	h := &Handler{svc: &service.Service{FAQ: faqSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/faqs")

	h.ListFAQs(c)

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

func TestFAQHandler_ListFAQs_Empty(t *testing.T) {
	mockRepo := &handlerMockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return nil, 0, nil
		},
	}

	faqSvc := service.NewFAQService(mockRepo)
	h := &Handler{svc: &service.Service{FAQ: faqSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/faqs")

	h.ListFAQs(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 for empty list, got %d", w.Code)
	}
}

func TestFAQHandler_ListFAQs_ServiceError(t *testing.T) {
	mockRepo := &handlerMockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return nil, 0, errors.New("db error")
		},
	}

	faqSvc := service.NewFAQService(mockRepo)
	h := &Handler{svc: &service.Service{FAQ: faqSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/faqs")

	h.ListFAQs(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500 for service error, got %d", w.Code)
	}
}

func TestFAQHandler_AdminListFAQs_Success(t *testing.T) {
	mockRepo := &handlerMockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return []repository.FAQWithProject{
				{FAQ: model.FAQ{ID: 1, Question: "Q1", Answer: "A1", CreatedAt: time.Now(), UpdatedAt: time.Now()}},
				{FAQ: model.FAQ{ID: 2, Question: "Q2", Answer: "A2", CreatedAt: time.Now(), UpdatedAt: time.Now()}},
			}, 2, nil
		},
	}

	faqSvc := service.NewFAQService(mockRepo)
	h := &Handler{svc: &service.Service{FAQ: faqSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/admin/faqs?page=1&per_page=10")

	h.AdminListFAQs(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var resp dto.PaginatedResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Code != 200 {
		t.Errorf("expected response code 200, got %d", resp.Code)
	}
}

func TestFAQHandler_AdminListFAQs_ServiceError(t *testing.T) {
	mockRepo := &handlerMockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return nil, 0, errors.New("db error")
		},
	}

	faqSvc := service.NewFAQService(mockRepo)
	h := &Handler{svc: &service.Service{FAQ: faqSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/admin/faqs?page=1&per_page=10")

	h.AdminListFAQs(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestFAQHandler_CreateFAQ_Success(t *testing.T) {
	mockRepo := &handlerMockFAQRepo{
		createFn: func(faq *model.FAQ) error {
			faq.ID = 100
			return nil
		},
	}

	faqSvc := service.NewFAQService(mockRepo)
	h := &Handler{svc: &service.Service{FAQ: faqSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makePostRequest("/api/v1/admin/faqs", model.FAQ{
		Question: "New question?",
		Answer:   "New answer.",
	})

	h.CreateFAQ(c)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestFAQHandler_CreateFAQ_InvalidJSON(t *testing.T) {
	h := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/admin/faqs", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateFAQ(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid JSON, got %d", w.Code)
	}
}

func TestFAQHandler_CreateFAQ_MissingFields(t *testing.T) {
	faqSvc := service.NewFAQService(&handlerMockFAQRepo{})
	h := &Handler{svc: &service.Service{FAQ: faqSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makePostRequest("/api/v1/admin/faqs", model.FAQ{})

	h.CreateFAQ(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for missing fields, got %d", w.Code)
	}
}

func TestFAQHandler_UpdateFAQ_Success(t *testing.T) {
	mockRepo := &handlerMockFAQRepo{
		findByIDFn: func(id uint64) (*model.FAQ, error) {
			return &model.FAQ{ID: 5, Question: "Old Q", Answer: "Old A"}, nil
		},
		updateFn: func(faq *model.FAQ) error {
			return nil
		},
	}

	faqSvc := service.NewFAQService(mockRepo)
	h := &Handler{svc: &service.Service{FAQ: faqSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makePostRequest("/api/v1/admin/faqs/5", model.FAQ{
		Question: "Updated Q",
		Answer:   "Updated A",
	})
	c.Params = gin.Params{{Key: "id", Value: "5"}}

	h.UpdateFAQ(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestFAQHandler_UpdateFAQ_InvalidID(t *testing.T) {
	h := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makePostRequest("/api/v1/admin/faqs/abc", model.FAQ{})
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h.UpdateFAQ(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid id, got %d", w.Code)
	}
}

func TestFAQHandler_DeleteFAQ_Success(t *testing.T) {
	mockRepo := &handlerMockFAQRepo{
		deleteFn: func(id uint64) error {
			return nil
		},
	}

	faqSvc := service.NewFAQService(mockRepo)
	h := &Handler{svc: &service.Service{FAQ: faqSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/admin/faqs/5")
	c.Params = gin.Params{{Key: "id", Value: "5"}}

	h.DeleteFAQ(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestFAQHandler_DeleteFAQ_InvalidID(t *testing.T) {
	h := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/admin/faqs/abc")
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h.DeleteFAQ(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid id, got %d", w.Code)
	}
}
