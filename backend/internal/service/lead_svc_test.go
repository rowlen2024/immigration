package service

import (
	"errors"
	"strings"
	"testing"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

// mockLeadRepo implements repository.LeadRepository for testing.
type mockLeadRepo struct {
	findAllFn      func(filter repository.LeadFilter) ([]model.Lead, int64, error)
	createFn       func(lead *model.Lead) error
	updateStatusFn func(id uint64, status string, notes string) error
	deleteFn       func(id uint64) error
}

func (m *mockLeadRepo) FindAll(filter repository.LeadFilter) ([]model.Lead, int64, error) {
	if m.findAllFn != nil {
		return m.findAllFn(filter)
	}
	return nil, 0, nil
}

func (m *mockLeadRepo) Create(lead *model.Lead) error {
	if m.createFn != nil {
		return m.createFn(lead)
	}
	return nil
}

func (m *mockLeadRepo) UpdateStatus(id uint64, status string, notes string) error {
	if m.updateStatusFn != nil {
		return m.updateStatusFn(id, status, notes)
	}
	return nil
}

func (m *mockLeadRepo) Delete(id uint64) error {
	if m.deleteFn != nil {
		return m.deleteFn(id)
	}
	return nil
}

func TestLead_Create_Success(t *testing.T) {
	var savedLead *model.Lead
	repo := &mockLeadRepo{
		createFn: func(lead *model.Lead) error {
			savedLead = lead
			lead.ID = 42
			return nil
		},
	}

	svc := NewLeadService(repo)

	lead, err := svc.Create(&dto.LeadRequest{
		Name:              "John Doe",
		Phone:             "1234567890",
		Email:             "john@example.com",
		InterestedProject: "project-a",
		Message:           "I am interested",
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if savedLead == nil {
		t.Fatal("expected Create to be called on repo")
	}
	if savedLead.Name != "John Doe" {
		t.Errorf("expected name 'John Doe', got '%s'", savedLead.Name)
	}
	if savedLead.Phone != "1234567890" {
		t.Errorf("expected phone '1234567890', got '%s'", savedLead.Phone)
	}
	if lead.ID != 42 {
		t.Errorf("expected ID 42, got %d", lead.ID)
	}
}

func TestLead_Create_XSSSanitization(t *testing.T) {
	var savedLead *model.Lead
	repo := &mockLeadRepo{
		createFn: func(lead *model.Lead) error {
			savedLead = lead
			return nil
		},
	}

	svc := NewLeadService(repo)

	_, err := svc.Create(&dto.LeadRequest{
		Name:    "<script>alert('xss')</script>John",
		Phone:   "1234567890",
		Message: "<b>Hello</b><script>bad</script>",
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if strings.Contains(savedLead.Name, "<script>") {
		t.Errorf("expected script tag in name to be sanitized")
	}
	if strings.Contains(savedLead.Message, "<script>") {
		t.Errorf("expected script tag in message to be sanitized")
	}
}

func TestLead_Create_NilRequest(t *testing.T) {
	repo := &mockLeadRepo{}
	svc := NewLeadService(repo)

	_, err := svc.Create(nil)
	if err == nil {
		t.Fatal("expected error for nil request")
	}
}

func TestLead_Create_MissingName(t *testing.T) {
	repo := &mockLeadRepo{}
	svc := NewLeadService(repo)

	_, err := svc.Create(&dto.LeadRequest{Phone: "1234567890"})
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestLead_Create_MissingPhone(t *testing.T) {
	repo := &mockLeadRepo{}
	svc := NewLeadService(repo)

	_, err := svc.Create(&dto.LeadRequest{Name: "John"})
	if err == nil {
		t.Fatal("expected error for missing phone")
	}
}

func TestLead_AdminList_Success(t *testing.T) {
	sampleLeads := []model.Lead{
		{ID: 1, Name: "Lead A", Status: "new"},
		{ID: 2, Name: "Lead B", Status: "contacted"},
	}

	repo := &mockLeadRepo{
		findAllFn: func(filter repository.LeadFilter) ([]model.Lead, int64, error) {
			return sampleLeads, int64(len(sampleLeads)), nil
		},
	}

	svc := NewLeadService(repo)

	leads, total, err := svc.List(dto.LeadListRequest{PaginationRequest: dto.PaginationRequest{Page: 1, PerPage: 10}})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(leads) != 2 {
		t.Errorf("expected 2 leads, got %d", len(leads))
	}
}

func TestLead_List_FilterByStatus(t *testing.T) {
	repo := &mockLeadRepo{
		findAllFn: func(filter repository.LeadFilter) ([]model.Lead, int64, error) {
			if filter.Status != "new" {
				t.Errorf("expected status filter 'new', got '%s'", filter.Status)
			}
			return []model.Lead{{ID: 1, Status: "new"}}, 1, nil
		},
	}

	svc := NewLeadService(repo)
	_, _, err := svc.List(dto.LeadListRequest{PaginationRequest: dto.PaginationRequest{Page: 1, PerPage: 10}, Status: "new"})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestLead_List_NoPagination(t *testing.T) {
	repo := &mockLeadRepo{
		findAllFn: func(filter repository.LeadFilter) ([]model.Lead, int64, error) {
			if filter.Page != 0 {
				t.Errorf("expected page 0 (no pagination), got %d", filter.Page)
			}
			return []model.Lead{}, 0, nil
		},
	}

	svc := NewLeadService(repo)
	_, _, err := svc.List(dto.LeadListRequest{})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestLead_UpdateStatus_Success(t *testing.T) {
	repo := &mockLeadRepo{
		updateStatusFn: func(id uint64, status string, notes string) error {
			return nil
		},
	}

	svc := NewLeadService(repo)

	lead, err := svc.Update(1, "contacted", "Called the customer")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if lead.ID != 1 {
		t.Errorf("expected ID 1, got %d", lead.ID)
	}
	if lead.Status != "contacted" {
		t.Errorf("expected status 'contacted', got '%s'", lead.Status)
	}
	if lead.Notes != "Called the customer" {
		t.Errorf("expected notes 'Called the customer', got '%s'", lead.Notes)
	}
}

func TestLead_UpdateStatus_InvalidStatus(t *testing.T) {
	repo := &mockLeadRepo{}
	svc := NewLeadService(repo)

	_, err := svc.Update(1, "invalid-status", "")
	if err == nil {
		t.Fatal("expected error for invalid status")
	}
}

func TestLead_UpdateStatus_ZeroID(t *testing.T) {
	repo := &mockLeadRepo{}
	svc := NewLeadService(repo)

	_, err := svc.Update(0, "new", "")
	if err == nil {
		t.Fatal("expected error for zero id")
	}
}

func TestLead_UpdateStatus_Qualified(t *testing.T) {
	repo := &mockLeadRepo{
		updateStatusFn: func(id uint64, status string, notes string) error {
			return nil
		},
	}

	svc := NewLeadService(repo)

	lead, err := svc.Update(10, "qualified", "Customer ready to proceed")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if lead.Status != "qualified" {
		t.Errorf("expected status 'qualified', got '%s'", lead.Status)
	}
}

func TestLead_UpdateStatus_Closed(t *testing.T) {
	repo := &mockLeadRepo{
		updateStatusFn: func(id uint64, status string, notes string) error {
			return nil
		},
	}

	svc := NewLeadService(repo)

	lead, err := svc.Update(10, "closed", "Deal closed")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if lead.Status != "closed" {
		t.Errorf("expected status 'closed', got '%s'", lead.Status)
	}
}

func TestLead_UpdateStatus_EmptyStatus(t *testing.T) {
	repo := &mockLeadRepo{}
	svc := NewLeadService(repo)

	_, err := svc.Update(1, "", "")
	if err == nil {
		t.Fatal("expected error for empty/invalid status")
	}
}

func TestLead_UpdateStatus_RepoError(t *testing.T) {
	repo := &mockLeadRepo{
		updateStatusFn: func(id uint64, status string, notes string) error {
			return errors.New("db error")
		},
	}

	svc := NewLeadService(repo)

	_, err := svc.Update(1, "contacted", "")
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestLead_List_RepoError(t *testing.T) {
	repo := &mockLeadRepo{
		findAllFn: func(filter repository.LeadFilter) ([]model.Lead, int64, error) {
			return nil, 0, errors.New("db error")
		},
	}

	svc := NewLeadService(repo)

	_, _, err := svc.List(dto.LeadListRequest{PaginationRequest: dto.PaginationRequest{Page: 1, PerPage: 10}})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestLead_Create_RepoError(t *testing.T) {
	repo := &mockLeadRepo{
		createFn: func(lead *model.Lead) error {
			return errors.New("db error")
		},
	}

	svc := NewLeadService(repo)

	_, err := svc.Create(&dto.LeadRequest{
		Name:  "John",
		Phone: "1234567890",
	})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestLead_Create_SpecialCharsInName(t *testing.T) {
	var savedLead *model.Lead
	repo := &mockLeadRepo{
		createFn: func(lead *model.Lead) error {
			savedLead = lead
			return nil
		},
	}

	svc := NewLeadService(repo)

	// Name with special HTML chars - strict policy should strip them
	_, err := svc.Create(&dto.LeadRequest{
		Name:    "John <b>Doe</b>",
		Phone:   "1234567890",
		Message: "Plain message",
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if savedLead == nil {
		t.Fatal("expected Create to be called on repo")
	}
	if strings.Contains(savedLead.Name, "<b>") {
		t.Errorf("expected HTML to be sanitized from name, got '%s'", savedLead.Name)
	}
}

func TestLead_Create_AllFieldsPopulated(t *testing.T) {
	var savedLead *model.Lead
	repo := &mockLeadRepo{
		createFn: func(lead *model.Lead) error {
			savedLead = lead
			lead.ID = 100
			return nil
		},
	}

	svc := NewLeadService(repo)

	lead, err := svc.Create(&dto.LeadRequest{
		Name:              "Alice Smith",
		Phone:             "9876543210",
		Email:             "alice@example.com",
		InterestedProject: "project-xyz",
		Message:           "I need more information about this project",
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if lead.ID != 100 {
		t.Errorf("expected ID 100, got %d", lead.ID)
	}
	if savedLead.Email != "alice@example.com" {
		t.Errorf("expected email 'alice@example.com', got '%s'", savedLead.Email)
	}
	if savedLead.InterestedProject != "project-xyz" {
		t.Errorf("expected project 'project-xyz', got '%s'", savedLead.InterestedProject)
	}
}
