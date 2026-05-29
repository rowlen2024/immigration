package service

import (
	"errors"
	"testing"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"time"
)

// mockProjectRepo implements repository.ProjectRepository for testing.
type mockProjectRepo struct {
	findByIDFn    func(id uint64) (*model.Project, error)
	findBySlugFn  func(slug string) (*model.Project, error)
	findAllFn     func(page, perPage int, search, status string) ([]model.Project, int64, error)
	findBySlugsFn func(slugs []string) ([]model.Project, error)
	createFn      func(project *model.Project) error
	updateFn      func(project *model.Project) error
	deleteFn      func(id uint64) error
	findNewsFn    func(projectID uint64) ([]model.Page, error)
	addNewsFn     func(projectID uint64, pageIDs []uint64) error
	removeNewsFn  func(projectID, pageID uint64) error
	deleteNewsByProjectIDFn func(projectID uint64) error
}

func (m *mockProjectRepo) FindByID(id uint64) (*model.Project, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return nil, errors.New("not found")
}

func (m *mockProjectRepo) FindBySlug(slug string) (*model.Project, error) {
	if m.findBySlugFn != nil {
		return m.findBySlugFn(slug)
	}
	return nil, errors.New("not found")
}

func (m *mockProjectRepo) FindAll(page, perPage int, search, status string) ([]model.Project, int64, error) {
	if m.findAllFn != nil {
		return m.findAllFn(page, perPage, search, status)
	}
	return nil, 0, nil
}

func (m *mockProjectRepo) FindAllWithoutPagination(search, status string) ([]model.Project, error) {
	return nil, nil
}

func (m *mockProjectRepo) FindBySlugs(slugs []string) ([]model.Project, error) {
	if m.findBySlugsFn != nil {
		return m.findBySlugsFn(slugs)
	}
	return nil, nil
}

func (m *mockProjectRepo) Create(project *model.Project) error {
	if m.createFn != nil {
		return m.createFn(project)
	}
	return nil
}

func (m *mockProjectRepo) Update(project *model.Project) error {
	if m.updateFn != nil {
		return m.updateFn(project)
	}
	return nil
}

func (m *mockProjectRepo) Delete(id uint64) error {
	if m.deleteFn != nil {
		return m.deleteFn(id)
	}
	return nil
}

func (m *mockProjectRepo) FindNews(projectID uint64) ([]model.Page, error) {
	if m.findNewsFn != nil {
		return m.findNewsFn(projectID)
	}
	return nil, nil
}

func (m *mockProjectRepo) AddNews(projectID uint64, pageIDs []uint64) error {
	if m.addNewsFn != nil {
		return m.addNewsFn(projectID, pageIDs)
	}
	return nil
}

func (m *mockProjectRepo) RemoveNews(projectID, pageID uint64) error {
	if m.removeNewsFn != nil {
		return m.removeNewsFn(projectID, pageID)
	}
	return nil
}

func (m *mockProjectRepo) DeleteNewsByProjectID(projectID uint64) error {
	if m.deleteNewsByProjectIDFn != nil {
		return m.deleteNewsByProjectIDFn(projectID)
	}
	return nil
}

func TestProject_List(t *testing.T) {
	sampleProjects := []model.Project{
		{ID: 1, Name: "Project A", Slug: "project-a"},
		{ID: 2, Name: "Project B", Slug: "project-b"},
		{ID: 3, Name: "Project C", Slug: "project-c"},
	}

	repo := &mockProjectRepo{
		findAllFn: func(page, perPage int, search, status string) ([]model.Project, int64, error) {
			return sampleProjects, int64(len(sampleProjects)), nil
		},
	}

	svc := NewProjectService(repo, nil)

	projects, total, err := svc.List(1, 10, "", "")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(projects) != 3 {
		t.Errorf("expected 3 projects, got %d", len(projects))
	}
}

func TestProject_List_DefaultPagination(t *testing.T) {
	repo := &mockProjectRepo{
		findAllFn: func(page, perPage int, search, status string) ([]model.Project, int64, error) {
			if page != 1 {
				t.Errorf("expected page 1 (default), got %d", page)
			}
			if perPage != 10 {
				t.Errorf("expected perPage 10 (default), got %d", perPage)
			}
			return []model.Project{}, 0, nil
		},
	}

	svc := NewProjectService(repo, nil)
	_, _, err := svc.List(0, 0, "", "")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestProject_GetBySlug_Success(t *testing.T) {
	expected := &model.Project{ID: 1, Name: "Test Project", Slug: "test-project"}
	repo := &mockProjectRepo{
		findBySlugFn: func(slug string) (*model.Project, error) {
			if slug == "test-project" {
				return expected, nil
			}
			return nil, errors.New("not found")
		},
	}

	svc := NewProjectService(repo, nil)

	project, err := svc.GetBySlug("test-project")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if project.ID != 1 || project.Name != "Test Project" {
		t.Errorf("unexpected project: %+v", project)
	}
}

func TestProject_GetBySlug_NotFound(t *testing.T) {
	repo := &mockProjectRepo{
		findBySlugFn: func(slug string) (*model.Project, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewProjectService(repo, nil)

	_, err := svc.GetBySlug("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent slug")
	}
}

func TestProject_GetBySlug_EmptySlug(t *testing.T) {
	repo := &mockProjectRepo{}
	svc := NewProjectService(repo, nil)

	_, err := svc.GetBySlug("")
	if err == nil {
		t.Fatal("expected error for empty slug")
	}
}

func TestProject_Compare_Success(t *testing.T) {
	expected := []model.Project{
		{ID: 1, Name: "Project A", Slug: "a"},
		{ID: 2, Name: "Project B", Slug: "b"},
	}

	repo := &mockProjectRepo{
		findBySlugsFn: func(slugs []string) ([]model.Project, error) {
			return expected, nil
		},
	}

	svc := NewProjectService(repo, nil)

	projects, err := svc.Compare([]string{"a", "b"})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(projects) != 2 {
		t.Errorf("expected 2 projects, got %d", len(projects))
	}
}

func TestProject_Compare_EmptySlugs(t *testing.T) {
	repo := &mockProjectRepo{}
	svc := NewProjectService(repo, nil)

	_, err := svc.Compare([]string{})
	if err == nil {
		t.Fatal("expected error for empty slugs")
	}
}

func TestProject_Compare_TooManySlugs(t *testing.T) {
	repo := &mockProjectRepo{}
	svc := NewProjectService(repo, nil)

	_, err := svc.Compare([]string{"a", "b", "c", "d", "e", "f"})
	if err == nil {
		t.Fatal("expected error for too many slugs")
	}
}

func TestProject_Create_Success(t *testing.T) {
	created := false
	repo := &mockProjectRepo{
		createFn: func(project *model.Project) error {
			created = true
			project.ID = 10
			return nil
		},
	}

	svc := NewProjectService(repo, nil)

	project, err := svc.Create(&model.Project{Name: "New Project", Slug: "new-project"})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !created {
		t.Error("expected Create to be called on repo")
	}
	if project.ID != 10 {
		t.Errorf("expected ID 10, got %d", project.ID)
	}
}

func TestProject_Create_NilProject(t *testing.T) {
	repo := &mockProjectRepo{}
	svc := NewProjectService(repo, nil)

	_, err := svc.Create(nil)
	if err == nil {
		t.Fatal("expected error for nil project")
	}
}

func TestProject_Create_MissingName(t *testing.T) {
	repo := &mockProjectRepo{}
	svc := NewProjectService(repo, nil)

	_, err := svc.Create(&model.Project{Slug: "some-slug"})
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestProject_Create_MissingSlug(t *testing.T) {
	repo := &mockProjectRepo{}
	svc := NewProjectService(repo, nil)

	_, err := svc.Create(&model.Project{Name: "Some Name"})
	if err == nil {
		t.Fatal("expected error for missing slug")
	}
}

func TestProject_Create_RepoError(t *testing.T) {
	repo := &mockProjectRepo{
		createFn: func(project *model.Project) error {
			return errors.New("db error")
		},
	}

	svc := NewProjectService(repo, nil)

	_, err := svc.Create(&model.Project{Name: "Test", Slug: "test"})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestProject_Update_Success(t *testing.T) {
	var updatedProject *model.Project
	repo := &mockProjectRepo{
		findByIDFn: func(id uint64) (*model.Project, error) {
			return &model.Project{ID: id, Slug: "existing-slug", Name: "Old Name"}, nil
		},
		updateFn: func(project *model.Project) error {
			updatedProject = project
			return nil
		},
	}

	svc := NewProjectService(repo, nil)

	project, err := svc.Update(5, dto.UpdateProjectRequest{Name: "Updated Name", Slug: "updated-slug"})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if project.ID != 5 {
		t.Errorf("expected ID 5, got %d", project.ID)
	}
	if updatedProject == nil {
		t.Fatal("expected Update to be called on repo")
	}
	if updatedProject.ID != 5 {
		t.Errorf("expected repo to receive ID 5, got %d", updatedProject.ID)
	}
}

func TestProject_Update_NilProject(t *testing.T) {
	repo := &mockProjectRepo{}
	svc := NewProjectService(repo, nil)

	_, err := svc.Update(1, dto.UpdateProjectRequest{})
	if err == nil {
		t.Fatal("expected error for nil project in update")
	}
}

func TestProject_Update_ZeroID(t *testing.T) {
	repo := &mockProjectRepo{}
	svc := NewProjectService(repo, nil)

	_, err := svc.Update(0, dto.UpdateProjectRequest{Name: "Test", Slug: "test"})
	if err == nil {
		t.Fatal("expected error for zero id")
	}
}

func TestProject_Update_RepoError(t *testing.T) {
	repo := &mockProjectRepo{
		updateFn: func(project *model.Project) error {
			return errors.New("db error")
		},
	}

	svc := NewProjectService(repo, nil)

	_, err := svc.Update(1, dto.UpdateProjectRequest{Name: "Test", Slug: "test"})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestProject_Delete_Success(t *testing.T) {
	deleted := false
	var deletedID uint64
	repo := &mockProjectRepo{
		deleteFn: func(id uint64) error {
			deleted = true
			deletedID = id
			return nil
		},
	}

	svc := NewProjectService(repo, nil)

	err := svc.Delete(10)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !deleted {
		t.Error("expected Delete to be called on repo")
	}
	if deletedID != 10 {
		t.Errorf("expected ID 10, got %d", deletedID)
	}
}

func TestProject_Delete_ZeroID(t *testing.T) {
	repo := &mockProjectRepo{}
	svc := NewProjectService(repo, nil)

	err := svc.Delete(0)
	if err == nil {
		t.Fatal("expected error for zero id")
	}
}

func TestProject_Delete_RepoError(t *testing.T) {
	repo := &mockProjectRepo{
		deleteFn: func(id uint64) error {
			return errors.New("db error")
		},
	}

	svc := NewProjectService(repo, nil)

	err := svc.Delete(1)
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestProject_AdminList_Success(t *testing.T) {
	sampleProjects := []model.Project{
		{ID: 1, Name: "Project A", Slug: "project-a"},
		{ID: 2, Name: "Project B", Slug: "project-b"},
	}

	repo := &mockProjectRepo{
		findAllFn: func(page, perPage int, search, status string) ([]model.Project, int64, error) {
			return sampleProjects, int64(len(sampleProjects)), nil
		},
	}

	svc := NewProjectService(repo, nil)

	projects, total, err := svc.AdminList(1, 10, "", "")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(projects) != 2 {
		t.Errorf("expected 2 projects, got %d", len(projects))
	}
}

func TestProject_AdminList_DefaultPagination(t *testing.T) {
	repo := &mockProjectRepo{
		findAllFn: func(page, perPage int, search, status string) ([]model.Project, int64, error) {
			if page != 1 {
				t.Errorf("expected page 1 (default), got %d", page)
			}
			if perPage != 10 {
				t.Errorf("expected perPage 10 (default), got %d", perPage)
			}
			return []model.Project{}, 0, nil
		},
	}

	svc := NewProjectService(repo, nil)
	_, _, err := svc.AdminList(0, 0, "", "")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestProject_List_RepoError(t *testing.T) {
	repo := &mockProjectRepo{
		findAllFn: func(page, perPage int, search, status string) ([]model.Project, int64, error) {
			return nil, 0, errors.New("db error")
		},
	}

	svc := NewProjectService(repo, nil)

	_, _, err := svc.List(1, 10, "", "")
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestProject_GetBySlug_RepoError(t *testing.T) {
	repo := &mockProjectRepo{
		findBySlugFn: func(slug string) (*model.Project, error) {
			return nil, errors.New("db error")
		},
	}

	svc := NewProjectService(repo, nil)

	_, err := svc.GetBySlug("valid-slug")
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestProject_Compare_RepoError(t *testing.T) {
	repo := &mockProjectRepo{
		findBySlugsFn: func(slugs []string) ([]model.Project, error) {
			return nil, errors.New("db error")
		},
	}

	svc := NewProjectService(repo, nil)

	_, err := svc.Compare([]string{"a", "b"})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func (m *mockProjectRepo) FindBySlugsLight(slugs []string) ([]model.Project, error) { return nil, nil }
func (m *mockProjectRepo) FindAllCoverImages() ([]string, error) { return nil, nil }
func (m *mockProjectRepo) Count() (int64, error) { return 0, nil }
func (m *mockProjectRepo) CountByRange(start, end time.Time) (int64, error) { return 0, nil }
