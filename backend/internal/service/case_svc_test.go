package service

import (
	"errors"
	"testing"

	"mygo-immigration/backend/internal/model"
)

// mockCaseRepo implements repository.CaseRepository for testing.
type mockCaseRepo struct {
	findByProjectIDFn func(projectID uint64) ([]model.Case, error)
	findAllFn         func(search string) ([]model.Case, error)
	createFn          func(c *model.Case) error
	updateFn          func(c *model.Case) error
	deleteFn          func(id uint64) error
	hardDeleteFn      func(id uint64) error
	findBySlugFn      func(slug string) (*model.Case, error)
}

func (m *mockCaseRepo) FindByProjectID(projectID uint64) ([]model.Case, error) {
	if m.findByProjectIDFn != nil {
		return m.findByProjectIDFn(projectID)
	}
	return nil, nil
}

func (m *mockCaseRepo) FindAll(search string) ([]model.Case, error) {
	if m.findAllFn != nil {
		return m.findAllFn(search)
	}
	return nil, nil
}

func (m *mockCaseRepo) Create(c *model.Case) error {
	if m.createFn != nil {
		return m.createFn(c)
	}
	return nil
}

func (m *mockCaseRepo) Update(c *model.Case) error {
	if m.updateFn != nil {
		return m.updateFn(c)
	}
	return nil
}

func (m *mockCaseRepo) Delete(id uint64) error {
	if m.deleteFn != nil {
		return m.deleteFn(id)
	}
	return nil
}

func (m *mockCaseRepo) HardDelete(id uint64) error {
	if m.hardDeleteFn != nil {
		return m.hardDeleteFn(id)
	}
	return nil
}

func (m *mockCaseRepo) FindBySlug(slug string) (*model.Case, error) {
	if m.findBySlugFn != nil {
		return m.findBySlugFn(slug)
	}
	return nil, nil
}

func TestCase_List(t *testing.T) {
	sampleCases := []model.Case{
		{ID: 1, Name: "Case A"},
		{ID: 2, Name: "Case B"},
		{ID: 3, Name: "Case C"},
	}

	repo := &mockCaseRepo{
		findAllFn: func(search string) ([]model.Case, error) {
			return sampleCases, nil
		},
	}

	svc := NewCaseService(repo)

	cases, err := svc.List()
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(cases) != 3 {
		t.Errorf("expected 3 cases, got %d", len(cases))
	}
}

func TestCase_List_Error(t *testing.T) {
	repo := &mockCaseRepo{
		findAllFn: func(search string) ([]model.Case, error) {
			return nil, errors.New("db error")
		},
	}

	svc := NewCaseService(repo)

	_, err := svc.List()
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestCase_Create_Success(t *testing.T) {
	created := false
	repo := &mockCaseRepo{
		createFn: func(c *model.Case) error {
			created = true
			c.ID = 50
			return nil
		},
	}

	svc := NewCaseService(repo)

	c, err := svc.Create(&model.Case{Name: "Test Case"})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !created {
		t.Error("expected Create to be called on repo")
	}
	if c.ID != 50 {
		t.Errorf("expected ID 50, got %d", c.ID)
	}
}

func TestCase_Create_NilCase(t *testing.T) {
	repo := &mockCaseRepo{}
	svc := NewCaseService(repo)

	_, err := svc.Create(nil)
	if err == nil {
		t.Fatal("expected error for nil case")
	}
}

func TestCase_Create_MissingName(t *testing.T) {
	repo := &mockCaseRepo{}
	svc := NewCaseService(repo)

	_, err := svc.Create(&model.Case{})
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestCase_Update_Success(t *testing.T) {
	updated := false
	repo := &mockCaseRepo{
		updateFn: func(c *model.Case) error {
			updated = true
			return nil
		},
	}

	svc := NewCaseService(repo)

	c, err := svc.Update(1, &model.Case{Name: "Updated Case"})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !updated {
		t.Error("expected Update to be called on repo")
	}
	if c.ID != 1 {
		t.Errorf("expected ID 1, got %d", c.ID)
	}
}

func TestCase_Update_NilCase(t *testing.T) {
	repo := &mockCaseRepo{}
	svc := NewCaseService(repo)

	_, err := svc.Update(1, nil)
	if err == nil {
		t.Fatal("expected error for nil case in update")
	}
}

func TestCase_Update_ZeroID(t *testing.T) {
	repo := &mockCaseRepo{}
	svc := NewCaseService(repo)

	_, err := svc.Update(0, &model.Case{Name: "Test"})
	if err == nil {
		t.Fatal("expected error for zero id")
	}
}

func TestCase_Update_RepoError(t *testing.T) {
	repo := &mockCaseRepo{
		updateFn: func(c *model.Case) error {
			return errors.New("db error")
		},
	}

	svc := NewCaseService(repo)

	_, err := svc.Update(1, &model.Case{Name: "Test"})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestCase_Delete_Success(t *testing.T) {
	deleted := false
	var deletedID uint64
	repo := &mockCaseRepo{
		deleteFn: func(id uint64) error {
			deleted = true
			deletedID = id
			return nil
		},
	}

	svc := NewCaseService(repo)

	err := svc.Delete(5)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !deleted {
		t.Error("expected Delete to be called on repo")
	}
	if deletedID != 5 {
		t.Errorf("expected ID 5, got %d", deletedID)
	}
}

func TestCase_Delete_ZeroID(t *testing.T) {
	repo := &mockCaseRepo{}
	svc := NewCaseService(repo)

	err := svc.Delete(0)
	if err == nil {
		t.Fatal("expected error for zero id")
	}
}

func TestCase_Delete_RepoError(t *testing.T) {
	repo := &mockCaseRepo{
		deleteFn: func(id uint64) error {
			return errors.New("db error")
		},
	}

	svc := NewCaseService(repo)

	err := svc.Delete(1)
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestCase_Create_RepoError(t *testing.T) {
	repo := &mockCaseRepo{
		createFn: func(c *model.Case) error {
			return errors.New("db error")
		},
	}

	svc := NewCaseService(repo)

	_, err := svc.Create(&model.Case{Name: "Test"})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestCase_AdminList_Success(t *testing.T) {
	sampleCases := []model.Case{
		{ID: 1, Name: "Case A"},
		{ID: 2, Name: "Case B"},
		{ID: 3, Name: "Case C"},
		{ID: 4, Name: "Case D"},
		{ID: 5, Name: "Case E"},
	}

	repo := &mockCaseRepo{
		findAllFn: func(search string) ([]model.Case, error) {
			return sampleCases, nil
		},
	}

	svc := NewCaseService(repo)

	cases, total, err := svc.AdminList(1, 2, "")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(cases) != 2 {
		t.Errorf("expected 2 cases on page 1 with perPage=2, got %d", len(cases))
	}
	if cases[0].ID != 1 || cases[1].ID != 2 {
		t.Errorf("expected first page to contain cases 1 and 2")
	}
}

func TestCase_AdminList_Page2(t *testing.T) {
	sampleCases := []model.Case{
		{ID: 1, Name: "Case A"},
		{ID: 2, Name: "Case B"},
		{ID: 3, Name: "Case C"},
	}

	repo := &mockCaseRepo{
		findAllFn: func(search string) ([]model.Case, error) {
			return sampleCases, nil
		},
	}

	svc := NewCaseService(repo)

	cases, total, err := svc.AdminList(2, 2, "")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(cases) != 1 {
		t.Errorf("expected 1 case on page 2 with perPage=2, got %d", len(cases))
	}
}

func TestCase_AdminList_BeyondRange(t *testing.T) {
	repo := &mockCaseRepo{
		findAllFn: func(search string) ([]model.Case, error) {
			return []model.Case{{ID: 1, Name: "Case A"}}, nil
		},
	}

	svc := NewCaseService(repo)

	cases, total, err := svc.AdminList(5, 10, "")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if len(cases) != 0 {
		t.Errorf("expected empty slice when page exceeds range, got %d cases", len(cases))
	}
}

func TestCase_AdminList_DefaultPagination(t *testing.T) {
	repo := &mockCaseRepo{
		findAllFn: func(search string) ([]model.Case, error) {
			return []model.Case{}, nil
		},
	}

	svc := NewCaseService(repo)
	_, _, err := svc.AdminList(0, 0, "")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestCase_AdminList_RepoError(t *testing.T) {
	repo := &mockCaseRepo{
		findAllFn: func(search string) ([]model.Case, error) {
			return nil, errors.New("db error")
		},
	}

	svc := NewCaseService(repo)

	_, _, err := svc.AdminList(1, 10, "")
	if err == nil {
		t.Fatal("expected error from repo")
	}
}
