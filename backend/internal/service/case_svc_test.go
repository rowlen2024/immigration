package service

import (
	"errors"
	"testing"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
	"time"
)

// mockCaseRepo implements repository.CaseRepository for testing.
type mockCaseRepo struct {
	findByIDFn          func(id uint64) (*model.Case, error)
	findBySlugFn        func(slug string) (*model.Case, error)
	findAllFn           func(filter repository.CaseFilter) ([]model.Case, int64, error)
	findOptionsFn       func(filter repository.CaseFilter) ([]repository.CaseOptionRow, int64, error)
	createFn            func(c *model.Case) error
	updateFn            func(c *model.Case) error
	deleteFn            func(id uint64) error
	deleteByProjectIDFn func(projectID uint64) error
}

func (m *mockCaseRepo) FindByID(id uint64) (*model.Case, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return nil, nil
}

func (m *mockCaseRepo) FindBySlug(slug string) (*model.Case, error) {
	if m.findBySlugFn != nil {
		return m.findBySlugFn(slug)
	}
	return nil, nil
}

func (m *mockCaseRepo) FindAll(filter repository.CaseFilter) ([]model.Case, int64, error) {
	if m.findAllFn != nil {
		return m.findAllFn(filter)
	}
	return nil, 0, nil
}

func (m *mockCaseRepo) FindOptions(filter repository.CaseFilter) ([]repository.CaseOptionRow, int64, error) {
	if m.findOptionsFn != nil {
		return m.findOptionsFn(filter)
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

func (m *mockCaseRepo) FindByIDs(ids []uint64) ([]model.Case, error)     { return nil, nil }
func (m *mockCaseRepo) FindAllPhotoURLs() ([]string, error)              { return nil, nil }
func (m *mockCaseRepo) FindAllContents() ([]string, error)               { return nil, nil }
func (m *mockCaseRepo) Count() (int64, error)                            { return 0, nil }
func (m *mockCaseRepo) CountByRange(start, end time.Time) (int64, error) { return 0, nil }

func TestCaseService_Options(t *testing.T) {
	projectID := uint64(9)
	repo := &mockCaseRepo{
		findOptionsFn: func(filter repository.CaseFilter) ([]repository.CaseOptionRow, int64, error) {
			if filter.ProjectID == nil || *filter.ProjectID != projectID {
				t.Fatalf("expected project_id %d, got %#v", projectID, filter.ProjectID)
			}
			if filter.Name != "case" {
				t.Fatalf("expected name filter case, got %q", filter.Name)
			}
			if filter.Page != 2 || filter.PerPage != 20 {
				t.Fatalf("expected pagination 2/20, got %d/%d", filter.Page, filter.PerPage)
			}
			return []repository.CaseOptionRow{{ID: 7, Name: "case one"}}, 1, nil
		},
	}
	svc := NewCaseService(repo, nil)

	items, total, err := svc.Options(dto.CaseOptionRequest{
		OptionPaginationRequest: dto.OptionPaginationRequest{Page: 2, PerPage: 20},
		ProjectID:               &projectID,
		Name:                    "case",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 1 || len(items) != 1 || items[0].ID != 7 {
		t.Fatalf("unexpected options result: total=%d items=%v", total, items)
	}
}

func TestCase_List(t *testing.T) {
	sampleCases := []model.Case{
		{ID: 1, Name: "Case A"},
		{ID: 2, Name: "Case B"},
		{ID: 3, Name: "Case C"},
	}

	repo := &mockCaseRepo{
		findAllFn: func(filter repository.CaseFilter) ([]model.Case, int64, error) {
			return sampleCases, 3, nil
		},
	}

	svc := NewCaseService(repo, nil)

	cases, total, err := svc.List(dto.CaseListRequest{})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(cases) != 3 {
		t.Errorf("expected 3 cases, got %d", len(cases))
	}
}

func TestCase_List_Error(t *testing.T) {
	repo := &mockCaseRepo{
		findAllFn: func(filter repository.CaseFilter) ([]model.Case, int64, error) {
			return nil, 0, errors.New("db error")
		},
	}

	svc := NewCaseService(repo, nil)

	_, _, err := svc.List(dto.CaseListRequest{})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestCase_List_Paginated(t *testing.T) {
	repo := &mockCaseRepo{
		findAllFn: func(filter repository.CaseFilter) ([]model.Case, int64, error) {
			return []model.Case{
				{ID: 1, Name: "Case A"},
				{ID: 2, Name: "Case B"},
			}, 5, nil
		},
	}

	svc := NewCaseService(repo, nil)

	cases, total, err := svc.List(dto.CaseListRequest{
		PaginationRequest: dto.PaginationRequest{Page: 1, PerPage: 2},
	})
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

func TestCase_List_Page2(t *testing.T) {
	repo := &mockCaseRepo{
		findAllFn: func(filter repository.CaseFilter) ([]model.Case, int64, error) {
			return []model.Case{
				{ID: 3, Name: "Case C"},
			}, 3, nil
		},
	}

	svc := NewCaseService(repo, nil)

	cases, total, err := svc.List(dto.CaseListRequest{
		PaginationRequest: dto.PaginationRequest{Page: 2, PerPage: 2},
	})
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

func TestCase_List_BeyondRange(t *testing.T) {
	repo := &mockCaseRepo{
		findAllFn: func(filter repository.CaseFilter) ([]model.Case, int64, error) {
			return []model.Case{}, 1, nil
		},
	}

	svc := NewCaseService(repo, nil)

	cases, total, err := svc.List(dto.CaseListRequest{
		PaginationRequest: dto.PaginationRequest{Page: 5, PerPage: 10},
	})
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

func TestCase_List_RepoError(t *testing.T) {
	repo := &mockCaseRepo{
		findAllFn: func(filter repository.CaseFilter) ([]model.Case, int64, error) {
			return nil, 0, errors.New("db error")
		},
	}

	svc := NewCaseService(repo, nil)

	_, _, err := svc.List(dto.CaseListRequest{
		PaginationRequest: dto.PaginationRequest{Page: 1, PerPage: 10},
	})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestCase_List_FilterByName(t *testing.T) {
	repo := &mockCaseRepo{
		findAllFn: func(filter repository.CaseFilter) ([]model.Case, int64, error) {
			if filter.Name != "test" {
				t.Errorf("expected Name='test', got '%s'", filter.Name)
			}
			return []model.Case{{ID: 1, Name: "Test Case"}}, 1, nil
		},
	}

	svc := NewCaseService(repo, nil)

	cases, total, err := svc.List(dto.CaseListRequest{Name: "test"})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if len(cases) != 1 {
		t.Errorf("expected 1 case, got %d", len(cases))
	}
}

func TestCase_List_FilterByProject(t *testing.T) {
	projectID := uint64(10)
	repo := &mockCaseRepo{
		findAllFn: func(filter repository.CaseFilter) ([]model.Case, int64, error) {
			if filter.ProjectID == nil || *filter.ProjectID != 10 {
				t.Errorf("expected ProjectID=10, got %v", filter.ProjectID)
			}
			return []model.Case{{ID: 1, Name: "Project Case"}}, 1, nil
		},
	}

	svc := NewCaseService(repo, nil)

	_, _, err := svc.List(dto.CaseListRequest{ProjectID: &projectID})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
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
