package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
	"mygo-immigration/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// handlerMockLeadRepo implements repository.LeadRepository.
type handlerMockLeadRepo struct {
	findAll      func(filter repository.LeadFilter) ([]model.Lead, int64, error)
	create       func(lead *model.Lead) error
	updateStatus func(id uint64, status string, notes string) error
	delete       func(id uint64) error
}

func (m *handlerMockLeadRepo) FindAll(filter repository.LeadFilter) ([]model.Lead, int64, error) {
	if m.findAll != nil {
		return m.findAll(filter)
	}
	return nil, 0, nil
}
func (m *handlerMockLeadRepo) Create(lead *model.Lead) error {
	if m.create != nil {
		return m.create(lead)
	}
	return nil
}
func (m *handlerMockLeadRepo) UpdateStatus(id uint64, status string, notes string) error {
	if m.updateStatus != nil {
		return m.updateStatus(id, status, notes)
	}
	return nil
}
func (m *handlerMockLeadRepo) Delete(id uint64) error {
	if m.delete != nil {
		return m.delete(id)
	}
	return nil
}

func TestLeadHandler_CreateLead_Success(t *testing.T) {
	mockRepo := &handlerMockLeadRepo{
		create: func(lead *model.Lead) error {
			lead.ID = 99
			return nil
		},
	}

	leadSvc := service.NewLeadService(mockRepo)
	h := &Handler{svc: &service.Service{Lead: leadSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makePostRequest("/api/v1/leads", dto.LeadRequest{
		Name:    "John Doe",
		Phone:   "1234567890",
		Email:   "john@example.com",
		Message: "Interested in immigration",
	})

	h.CreateLead(c)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d, body: %s", w.Code, w.Body.String())
	}

	var resp dto.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Code != 200 {
		t.Errorf("expected response code 200, got %d", resp.Code)
	}
}

func TestLeadHandler_CreateLead_InvalidJSON(t *testing.T) {
	h := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/leads", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateLead(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid JSON, got %d", w.Code)
	}
}

func TestLeadHandler_CreateLead_MissingRequired(t *testing.T) {
	leadSvc := service.NewLeadService(&handlerMockLeadRepo{})
	h := &Handler{svc: &service.Service{Lead: leadSvc}}

	// Send a lead missing phone - the dto.LeadRequest has Phone with `binding:"required"`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Use ShouldBindJSON - it will fail because phone is required
	c.Request = makePostRequest("/api/v1/leads", map[string]string{
		"name": "",
	})

	h.CreateLead(c)

	// The handler will get error from ShouldBindJSON and return 400
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for missing required fields, got %d, body: %s", w.Code, w.Body.String())
	}
}
