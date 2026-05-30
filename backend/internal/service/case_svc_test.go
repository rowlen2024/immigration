package service

import (
	"errors"
	"testing"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"time"
)

// mockCaseRepo implements repository.CaseRepository for testing.
type mockCaseRepo struct {
	findByIDFn        func(id uint64) (*model.Case, error)
	findByProjectIDFn func(projectID uint64) ([]model.Case, error)
	findAllFn           func(search string) ([]model.Case, error)
	findAllPaginatedFn  func(page, perPage int, search string) ([]model.Case, int64, error)
	findFilteredPaginatedFn func(projectID *uint64, countryFrom string, page, perPage int) ([]model.Case, int64, error)
	createFn          func(c *model.Case) error
	updateFn          func(c *model.Case) error
	deleteFn          func(id uint64) error
	hardDeleteFn      func(id uint64) error
	deleteByProjectIDFn func(projectID uint64) error
	findBySlugFn      func(slug string) (*model.Case, error)
}

func (m *mockCaseRepo) FindByID(id uint64) (*model.Case, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return nil, nil
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

func (m *mockCaseRepo) FindAllPaginated(page, perPage int, search string) ([]model.Case, int64, error) {
	if m.findAllPaginatedFn != nil {
		return m.findAllPaginatedFn(page, perPage, search)
	}
	return nil, 0, nil
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

func (m *mockCaseRepo) DeleteByProjectID(projectID uint64) error {
	if m.deleteByProjectIDFn != nil {
		return m.deleteByProjectIDFn(projectID)
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

func (m *mockCaseRepo) FindFilteredPaginated(projectID *uint64, countryFrom string, page, perPage int) ([]model.Case, int64, error) {
	if m.findFilteredPaginatedFn != nil {
		return m.findFilteredPaginatedFn(projectID, countryFrom, page, perPage)
	}
	return nil, 0, nil
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

	svc := NewCaseService(repo, nil)

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

	svc := NewCaseService(repo, nil)

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

	svc := NewCaseService(repo, nil)

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
	svc := NewCaseService(repo, nil)

	_, err := svc.Create(nil)
	if err == nil {
		t.Fatal("expected error for nil case")
	}
}

func TestCase_Create_MissingName(t *testing.T) {
	repo := &mockCaseRepo{}
	svc := NewCaseService(repo, nil)

	_, err := svc.Create(&model.Case{})
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestCase_Update_Success(t *testing.T) {
	updated := false
	repo := &mockCaseRepo{
		findByIDFn: func(id uint64) (*model.Case, error) {
			return &model.Case{ID: id, Slug: "existing-slug"}, nil
		},
		updateFn: func(c *model.Case) error {
			updated = true
			return nil
		},
	}

	svc := NewCaseService(repo, nil)

	c, err := svc.Update(1, dto.UpdateCaseRequest{Name: "Updated Case"})
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

func TestCase_Update_NotFound(t *testing.T) {
	repo := &mockCaseRepo{
		findByIDFn: func(id uint64) (*model.Case, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewCaseService(repo, nil)

	_, err := svc.Update(1, dto.UpdateCaseRequest{Name: "Test"})
	if err == nil {
		t.Fatal("expected error for non-existent case")
	}
}

func TestCase_Update_ZeroID(t *testing.T) {
	repo := &mockCaseRepo{}
	svc := NewCaseService(repo, nil)

	_, err := svc.Update(0, dto.UpdateCaseRequest{Name: "Test"})
	if err == nil {
		t.Fatal("expected error for zero id")
	}
}

func TestCase_Update_RepoError(t *testing.T) {
	repo := &mockCaseRepo{
		findByIDFn: func(id uint64) (*model.Case, error) {
			return &model.Case{ID: id, Slug: "existing"}, nil
		},
		updateFn: func(c *model.Case) error {
			return errors.New("db error")
		},
	}

	svc := NewCaseService(repo, nil)

	_, err := svc.Update(1, dto.UpdateCaseRequest{Name: "Test"})
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

	svc := NewCaseService(repo, nil)

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
	svc := NewCaseService(repo, nil)

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

	svc := NewCaseService(repo, nil)

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

	svc := NewCaseService(repo, nil)

	_, err := svc.Create(&model.Case{Name: "Test"})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestCase_AdminList_Success(t *testing.T) {
	repo := &mockCaseRepo{
		findAllPaginatedFn: func(page, perPage int, search string) ([]model.Case, int64, error) {
			return []model.Case{
				{ID: 1, Name: "Case A"},
				{ID: 2, Name: "Case B"},
			}, 5, nil
		},
	}

	svc := NewCaseService(repo, nil)

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
}

func TestCase_AdminList_Page2(t *testing.T) {
	repo := &mockCaseRepo{
		findAllPaginatedFn: func(page, perPage int, search string) ([]model.Case, int64, error) {
			return []model.Case{
				{ID: 3, Name: "Case C"},
			}, 3, nil
		},
	}

	svc := NewCaseService(repo, nil)

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
		findAllPaginatedFn: func(page, perPage int, search string) ([]model.Case, int64, error) {
			return []model.Case{}, 1, nil
		},
	}

	svc := NewCaseService(repo, nil)

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
		findAllPaginatedFn: func(page, perPage int, search string) ([]model.Case, int64, error) {
			return []model.Case{}, 0, nil
		},
	}

	svc := NewCaseService(repo, nil)
	_, _, err := svc.AdminList(0, 0, "")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestCase_AdminList_RepoError(t *testing.T) {
	repo := &mockCaseRepo{
		findAllPaginatedFn: func(page, perPage int, search string) ([]model.Case, int64, error) {
			return nil, 0, errors.New("db error")
		},
	}

	svc := NewCaseService(repo, nil)

	_, _, err := svc.AdminList(1, 10, "")
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestCase_ListPaginated_Success(t *testing.T) {
	repo := &mockCaseRepo{
		findAllPaginatedFn: func(page, perPage int, search string) ([]model.Case, int64, error) {
			return []model.Case{
				{ID: 1, Name: "Case A"},
				{ID: 2, Name: "Case B"},
			}, 2, nil
		},
	}

	svc := NewCaseService(repo, nil)

	cases, total, err := svc.ListPaginated(1, 10)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(cases) != 2 {
		t.Errorf("expected 2 cases, got %d", len(cases))
	}
}

func (m *mockCaseRepo) FindByIDs(ids []uint64) ([]model.Case, error) { return nil, nil }
func (m *mockCaseRepo) FindAllPhotoURLs() ([]string, error) { return nil, nil }
func (m *mockCaseRepo) FindAllContents() ([]string, error) { return nil, nil }
func (m *mockCaseRepo) Count() (int64, error) { return 0, nil }
func (m *mockCaseRepo) CountByRange(start, end time.Time) (int64, error) { return 0, nil }
