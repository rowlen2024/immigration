package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
	"mygo-immigration/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// userHandlerMockUserRepo implements repository.UserRepository.
type userHandlerMockUserRepo struct {
	findByUsername func(username string) (*model.User, error)
	findAll        func(filter repository.UserFilter) ([]model.User, int64, error)
	create         func(user *model.User) error
	update         func(user *model.User) error
	findByID       func(id uint64) (*model.User, error)
	patchUpdate    func(id uint64, updates map[string]interface{}) error
	delete         func(id uint64) error
}

func (m *userHandlerMockUserRepo) FindByUsername(username string) (*model.User, error) {
	if m.findByUsername != nil {
		return m.findByUsername(username)
	}
	return nil, nil
}
func (m *userHandlerMockUserRepo) FindAll(filter repository.UserFilter) ([]model.User, int64, error) {
	if m.findAll != nil {
		return m.findAll(filter)
	}
	return nil, 0, nil
}
func (m *userHandlerMockUserRepo) Create(user *model.User) error {
	if m.create != nil {
		return m.create(user)
	}
	return nil
}
func (m *userHandlerMockUserRepo) Update(user *model.User) error {
	if m.update != nil {
		return m.update(user)
	}
	return nil
}
func (m *userHandlerMockUserRepo) FindByID(id uint64) (*model.User, error) {
	if m.findByID != nil {
		return m.findByID(id)
	}
	return nil, nil
}
func (m *userHandlerMockUserRepo) PatchUpdate(id uint64, updates map[string]interface{}) error {
	if m.patchUpdate != nil {
		return m.patchUpdate(id, updates)
	}
	return nil
}
func (m *userHandlerMockUserRepo) Delete(id uint64) error {
	if m.delete != nil {
		return m.delete(id)
	}
	return nil
}

func TestUserHandler_AdminListUsers_Success(t *testing.T) {
	mockRepo := &userHandlerMockUserRepo{
		findAll: func(filter repository.UserFilter) ([]model.User, int64, error) {
			return []model.User{
				{ID: 1, Username: "admin", Role: "admin"},
				{ID: 2, Username: "editor", Role: "editor"},
			}, 2, nil
		},
	}

	userSvc := service.NewUserService(mockRepo)
	h := &Handler{svc: &service.Service{User: userSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/admin/users?page=1&per_page=10")

	h.AdminListUsers(c)

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

func TestUserHandler_AdminGetUser_Success(t *testing.T) {
	mockRepo := &userHandlerMockUserRepo{
		findByID: func(id uint64) (*model.User, error) {
			return &model.User{ID: id, Username: "admin", Role: "admin"}, nil
		},
	}

	userSvc := service.NewUserService(mockRepo)
	h := &Handler{svc: &service.Service{User: userSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/admin/users/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.AdminGetUser(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestUserHandler_AdminGetUser_InvalidID(t *testing.T) {
	mockRepo := &userHandlerMockUserRepo{}
	userSvc := service.NewUserService(mockRepo)
	h := &Handler{svc: &service.Service{User: userSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/admin/users/abc", nil)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h.AdminGetUser(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestUserHandler_AdminCreateUser_Success(t *testing.T) {
	mockRepo := &userHandlerMockUserRepo{
		create: func(user *model.User) error {
			user.ID = 1
			return nil
		},
	}

	userSvc := service.NewUserService(mockRepo)
	h := &Handler{svc: &service.Service{User: userSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makePostRequest("/api/v1/admin/users", dto.CreateUserRequest{
		Username:    "newuser",
		Password:    "password123",
		DisplayName: "New User",
		Role:        "editor",
	})

	h.AdminCreateUser(c)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUserHandler_AdminCreateUser_InvalidJSON(t *testing.T) {
	h := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/admin/users", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	h.AdminCreateUser(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid JSON, got %d", w.Code)
	}
}

func TestUserHandler_AdminUpdateUser_Success(t *testing.T) {
	mockRepo := &userHandlerMockUserRepo{
		findByID: func(id uint64) (*model.User, error) {
			return &model.User{ID: id, Username: "existing"}, nil
		},
		update: func(user *model.User) error {
			return nil
		},
	}

	userSvc := service.NewUserService(mockRepo)
	h := &Handler{svc: &service.Service{User: userSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makePostRequest("/api/v1/admin/users/1", dto.UpdateUserRequest{
		DisplayName: "Updated",
		Role:        "admin",
	})
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.AdminUpdateUser(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUserHandler_AdminUpdateUser_InvalidID(t *testing.T) {
	mockRepo := &userHandlerMockUserRepo{}
	userSvc := service.NewUserService(mockRepo)
	h := &Handler{svc: &service.Service{User: userSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/admin/users/abc", nil)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h.AdminUpdateUser(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestUserHandler_AdminChangeMyPassword_Success(t *testing.T) {
	hashed := mustHashPassword("oldpass")
	mockRepo := &userHandlerMockUserRepo{
		findByID: func(id uint64) (*model.User, error) {
			return &model.User{ID: id, Username: "user", PasswordHash: hashed}, nil
		},
		update: func(user *model.User) error {
			return nil
		},
	}

	userSvc := service.NewUserService(mockRepo)
	h := &Handler{svc: &service.Service{User: userSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", uint64(1))
	c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/admin/me/password", strings.NewReader(`{"old_password":"oldpass","new_password":"newpass"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.AdminChangeMyPassword(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUserHandler_AdminChangeMyPassword_Unauthorized(t *testing.T) {
	userSvc := service.NewUserService(&userHandlerMockUserRepo{})
	h := &Handler{svc: &service.Service{User: userSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/admin/me/password", strings.NewReader(`{"old_password":"oldpass","new_password":"newpass"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.AdminChangeMyPassword(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestUserHandler_AdminChangeMyPassword_InvalidRequest(t *testing.T) {
	userSvc := service.NewUserService(&userHandlerMockUserRepo{})
	h := &Handler{svc: &service.Service{User: userSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", uint64(1))
	c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/admin/me/password", strings.NewReader(`{"old_password":"oldpass","new_password":"123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.AdminChangeMyPassword(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestUserHandler_AdminDeleteUser_Success(t *testing.T) {
	mockRepo := &userHandlerMockUserRepo{
		delete: func(id uint64) error {
			return nil
		},
	}

	userSvc := service.NewUserService(mockRepo)
	h := &Handler{svc: &service.Service{User: userSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/admin/users/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.AdminDeleteUser(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestUserHandler_AdminDeleteUser_InvalidID(t *testing.T) {
	mockRepo := &userHandlerMockUserRepo{}
	userSvc := service.NewUserService(mockRepo)
	h := &Handler{svc: &service.Service{User: userSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/admin/users/abc", nil)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h.AdminDeleteUser(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}
