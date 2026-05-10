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

// handlerMockProjectRepo implements repository.ProjectRepository.
type handlerMockProjectRepo struct {
	findBySlug  func(slug string) (*model.Project, error)
	findAll     func(page, perPage int, search, status string) ([]model.Project, int64, error)
	findBySlugs func(slugs []string) ([]model.Project, error)
	create      func(project *model.Project) error
	update      func(project *model.Project) error
	delete      func(id uint64) error
	findNews    func(projectID uint64) ([]model.Page, error)
	addNews     func(projectID uint64, pageIDs []uint64) error
	removeNews  func(projectID, pageID uint64) error
}

func (m *handlerMockProjectRepo) FindBySlug(slug string) (*model.Project, error) {
	if m.findBySlug != nil {
		return m.findBySlug(slug)
	}
	return nil, errors.New("not found")
}
func (m *handlerMockProjectRepo) FindAll(page, perPage int, search, status string) ([]model.Project, int64, error) {
	if m.findAll != nil {
		return m.findAll(page, perPage, search, status)
	}
	return nil, 0, nil
}
func (m *handlerMockProjectRepo) FindBySlugs(slugs []string) ([]model.Project, error) {
	if m.findBySlugs != nil {
		return m.findBySlugs(slugs)
	}
	return nil, nil
}
func (m *handlerMockProjectRepo) Create(project *model.Project) error {
	if m.create != nil {
		return m.create(project)
	}
	return nil
}
func (m *handlerMockProjectRepo) Update(project *model.Project) error {
	if m.update != nil {
		return m.update(project)
	}
	return nil
}
func (m *handlerMockProjectRepo) Delete(id uint64) error {
	if m.delete != nil {
		return m.delete(id)
	}
	return nil
}
func (m *handlerMockProjectRepo) FindNews(projectID uint64) ([]model.Page, error) {
	if m.findNews != nil {
		return m.findNews(projectID)
	}
	return nil, nil
}
func (m *handlerMockProjectRepo) AddNews(projectID uint64, pageIDs []uint64) error {
	if m.addNews != nil {
		return m.addNews(projectID, pageIDs)
	}
	return nil
}
func (m *handlerMockProjectRepo) RemoveNews(projectID, pageID uint64) error {
	if m.removeNews != nil {
		return m.removeNews(projectID, pageID)
	}
	return nil
}

func TestProjectHandler_ListProjects(t *testing.T) {
	mockRepo := &handlerMockProjectRepo{
		findAll: func(page, perPage int, search, status string) ([]model.Project, int64, error) {
			return []model.Project{
				{ID: 1, Name: "Project A", Slug: "project-a"},
				{ID: 2, Name: "Project B", Slug: "project-b"},
			}, 2, nil
		},
	}

	projectSvc := service.NewProjectService(mockRepo)
	h := &Handler{svc: &service.Service{Project: projectSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/projects?page=1&per_page=10")

	h.ListProjects(c)

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

func TestProjectHandler_GetProject_Success(t *testing.T) {
	mockRepo := &handlerMockProjectRepo{
		findBySlug: func(slug string) (*model.Project, error) {
			return &model.Project{ID: 1, Name: "Test Project", Slug: slug}, nil
		},
	}

	projectSvc := service.NewProjectService(mockRepo)
	h := &Handler{svc: &service.Service{Project: projectSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/projects/test-project")
	c.Params = gin.Params{{Key: "slug", Value: "test-project"}}

	h.GetProject(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var resp dto.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
}

func TestProjectHandler_GetProject_NotFound(t *testing.T) {
	mockRepo := &handlerMockProjectRepo{
		findBySlug: func(slug string) (*model.Project, error) {
			return nil, errors.New("record not found")
		},
	}

	projectSvc := service.NewProjectService(mockRepo)
	h := &Handler{svc: &service.Service{Project: projectSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/projects/nonexistent")
	c.Params = gin.Params{{Key: "slug", Value: "nonexistent"}}

	h.GetProject(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestProjectHandler_GetProject_MissingSlug(t *testing.T) {
	h := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/projects/")
	// No params set - slug is empty

	h.GetProject(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestProjectHandler_CompareProjects_Success(t *testing.T) {
	mockRepo := &handlerMockProjectRepo{
		findBySlugs: func(slugs []string) ([]model.Project, error) {
			return []model.Project{
				{ID: 1, Name: "Project A", Slug: "a"},
				{ID: 2, Name: "Project B", Slug: "b"},
			}, nil
		},
	}

	projectSvc := service.NewProjectService(mockRepo)
	h := &Handler{svc: &service.Service{Project: projectSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/projects/compare?slugs=a,b")

	h.CompareProjects(c)

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

func TestProjectHandler_CompareProjects_NoSlugs(t *testing.T) {
	h := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/projects/compare")

	h.CompareProjects(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for missing slugs, got %d", w.Code)
	}
}

func TestProjectHandler_CompareProjects_TooMany(t *testing.T) {
	mockRepo := &handlerMockProjectRepo{}
	projectSvc := service.NewProjectService(mockRepo)
	h := &Handler{svc: &service.Service{Project: projectSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/projects/compare?slugs=a,b,c,d,e,f")

	h.CompareProjects(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500 for too many slugs, got %d", w.Code)
	}
}

func TestProjectHandler_CompareProjects_ThreeWay(t *testing.T) {
	mockRepo := &handlerMockProjectRepo{
		findBySlugs: func(slugs []string) ([]model.Project, error) {
			return []model.Project{
				{ID: 1, Name: "Project A", Slug: "a", InvestmentAmount: "100万", ProcessingPeriod: "12月", TargetCrowd: "投资者"},
				{ID: 2, Name: "Project B", Slug: "b", InvestmentAmount: "200万", ProcessingPeriod: "24月", TargetCrowd: "企业家"},
				{ID: 3, Name: "Project C", Slug: "c", InvestmentAmount: "300万", ProcessingPeriod: "36月", TargetCrowd: "高净值"},
			}, nil
		},
	}
	projectSvc := service.NewProjectService(mockRepo)
	h := &Handler{svc: &service.Service{Project: projectSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/projects/compare?slugs=a,b,c")

	h.CompareProjects(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	data := body["data"].(map[string]interface{})
	projects := data["projects"].([]interface{})
	if len(projects) != 3 {
		t.Errorf("expected 3 projects in response, got %d", len(projects))
	}
	rows := data["rows"].([]interface{})
	firstRow := rows[0].(map[string]interface{})
	values := firstRow["values"].([]interface{})
	if len(values) != 3 {
		t.Errorf("expected 3 values per row, got %d", len(values))
	}
}

func TestProjectHandler_AdminListProjects_Success(t *testing.T) {
	mockRepo := &handlerMockProjectRepo{
		findAll: func(page, perPage int, search, status string) ([]model.Project, int64, error) {
			return []model.Project{
				{ID: 1, Name: "Project A", Slug: "project-a"},
				{ID: 2, Name: "Project B", Slug: "project-b"},
			}, 2, nil
		},
	}

	projectSvc := service.NewProjectService(mockRepo)
	h := &Handler{svc: &service.Service{Project: projectSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/admin/projects?page=1&per_page=10")

	h.AdminListProjects(c)

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

func TestProjectHandler_CreateProject_Success(t *testing.T) {
	mockRepo := &handlerMockProjectRepo{
		create: func(project *model.Project) error {
			project.ID = 99
			return nil
		},
	}

	projectSvc := service.NewProjectService(mockRepo)
	h := &Handler{svc: &service.Service{Project: projectSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makePostRequest("/api/v1/admin/projects", model.Project{
		Name: "New Test Project",
		Slug: "new-test-project",
	})

	h.CreateProject(c)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestProjectHandler_CreateProject_InvalidJSON(t *testing.T) {
	h := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/admin/projects", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateProject(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid JSON, got %d", w.Code)
	}
}

func TestProjectHandler_CreateProject_MissingFields(t *testing.T) {
	mockRepo := &handlerMockProjectRepo{}
	projectSvc := service.NewProjectService(mockRepo)
	h := &Handler{svc: &service.Service{Project: projectSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Send project with empty name and slug
	c.Request = makePostRequest("/api/v1/admin/projects", model.Project{
		Name: "",
		Slug: "",
	})

	h.CreateProject(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500 for missing fields, got %d", w.Code)
	}
}

func TestProjectHandler_UpdateProject_InvalidID(t *testing.T) {
	h := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makePostRequest("/api/v1/admin/projects/abc", model.Project{})
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h.UpdateProject(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid id, got %d", w.Code)
	}
}

func TestProjectHandler_UpdateProject_Success(t *testing.T) {
	mockRepo := &handlerMockProjectRepo{
		update: func(project *model.Project) error {
			return nil
		},
	}

	projectSvc := service.NewProjectService(mockRepo)
	h := &Handler{svc: &service.Service{Project: projectSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makePostRequest("/api/v1/admin/projects/5", model.Project{
		Name: "Updated Name",
		Slug: "updated-slug",
	})
	c.Params = gin.Params{{Key: "id", Value: "5"}}

	h.UpdateProject(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestProjectHandler_DeleteProject_InvalidID(t *testing.T) {
	h := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/admin/projects/abc")
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h.DeleteProject(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid id, got %d", w.Code)
	}
}

func TestProjectHandler_DeleteProject_Success(t *testing.T) {
	mockRepo := &handlerMockProjectRepo{
		delete: func(id uint64) error {
			return nil
		},
	}

	projectSvc := service.NewProjectService(mockRepo)
	h := &Handler{svc: &service.Service{Project: projectSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/admin/projects/5")
	c.Params = gin.Params{{Key: "id", Value: "5"}}

	h.DeleteProject(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestProjectHandler_ListProjects_ServiceError(t *testing.T) {
	mockRepo := &handlerMockProjectRepo{
		findAll: func(page, perPage int, search, status string) ([]model.Project, int64, error) {
			return nil, 0, errors.New("db error")
		},
	}

	projectSvc := service.NewProjectService(mockRepo)
	h := &Handler{svc: &service.Service{Project: projectSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/projects?page=1&per_page=10")

	h.ListProjects(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d, body: %s", w.Code, w.Body.String())
	}
}
