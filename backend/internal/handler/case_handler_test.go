package handler

import (
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
)

type handlerMockCaseRepo struct {
	findByIDFn func(id uint64) (*model.Case, error)
	createFn   func(c *model.Case) error
	updateFn   func(c *model.Case) error
}

func (m *handlerMockCaseRepo) FindByID(id uint64) (*model.Case, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return nil, nil
}
func (m *handlerMockCaseRepo) FindBySlug(string) (*model.Case, error)   { return nil, nil }
func (m *handlerMockCaseRepo) FindByIDs([]uint64) ([]model.Case, error) { return nil, nil }
func (m *handlerMockCaseRepo) FindAll(repository.CaseFilter) ([]model.Case, int64, error) {
	return nil, 0, nil
}
func (m *handlerMockCaseRepo) FindOptions(repository.CaseFilter) ([]repository.CaseOptionRow, int64, error) {
	return nil, 0, nil
}
func (m *handlerMockCaseRepo) Create(c *model.Case) error {
	if m.createFn != nil {
		return m.createFn(c)
	}
	return nil
}
func (m *handlerMockCaseRepo) Update(c *model.Case) error {
	if m.updateFn != nil {
		return m.updateFn(c)
	}
	return nil
}
func (m *handlerMockCaseRepo) Delete(uint64) error                 { return nil }
func (m *handlerMockCaseRepo) DeleteByProjectID(uint64) error      { return nil }
func (m *handlerMockCaseRepo) FindAllPhotoURLs() ([]string, error) { return nil, nil }
func (m *handlerMockCaseRepo) FindAllContents() ([]string, error)  { return nil, nil }
func (m *handlerMockCaseRepo) Count() (int64, error)               { return 0, nil }
func (m *handlerMockCaseRepo) CountByRange(time.Time, time.Time) (int64, error) {
	return 0, nil
}

func TestCaseHandlerCreateBindsPinned(t *testing.T) {
	var createdPinned bool
	repo := &handlerMockCaseRepo{createFn: func(c *model.Case) error {
		createdPinned = c.IsPinned
		c.ID = 99
		return nil
	}}
	h := &Handler{svc: &service.Service{Case: service.NewCaseService(repo, nil)}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makePostRequest("/api/v1/admin/cases", model.Case{Name: "New Case", IsPinned: true})

	h.CreateCase(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d, body: %s", w.Code, w.Body.String())
	}
	if !createdPinned {
		t.Error("expected create handler to bind is_pinned")
	}
	if !strings.Contains(w.Body.String(), `"is_pinned":true`) {
		t.Errorf("expected response to include is_pinned, body: %s", w.Body.String())
	}
}

func TestCaseHandlerUpdateBindsPinned(t *testing.T) {
	var updatedPinned bool
	repo := &handlerMockCaseRepo{
		findByIDFn: func(id uint64) (*model.Case, error) {
			return &model.Case{ID: id, Name: "Old Case"}, nil
		},
		updateFn: func(c *model.Case) error {
			updatedPinned = c.IsPinned
			return nil
		},
	}
	h := &Handler{svc: &service.Service{Case: service.NewCaseService(repo, nil)}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makePostRequest("/api/v1/admin/cases/5", dto.UpdateCaseRequest{Name: "Updated Case", IsPinned: true})
	c.Params = gin.Params{{Key: "id", Value: "5"}}

	h.UpdateCase(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
	if !updatedPinned {
		t.Error("expected update handler to bind is_pinned")
	}
	if !strings.Contains(w.Body.String(), `"is_pinned":true`) {
		t.Errorf("expected response to include is_pinned, body: %s", w.Body.String())
	}
}
