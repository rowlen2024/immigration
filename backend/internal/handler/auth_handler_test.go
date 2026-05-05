package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"mygo-immigration/backend/internal/config"
	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// handlerMockUserRepo implements repository.UserRepository.
type handlerMockUserRepo struct {
	findByUsername func(username string) (*model.User, error)
}

func (m *handlerMockUserRepo) FindByUsername(username string) (*model.User, error) {
	return m.findByUsername(username)
}

func (m *handlerMockUserRepo) FindAll() ([]model.User, error)                    { return nil, nil }
func (m *handlerMockUserRepo) FindAllPaginated(page, perPage int) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (m *handlerMockUserRepo) Create(user *model.User) error                     { return nil }
func (m *handlerMockUserRepo) Update(user *model.User) error                     { return nil }
func (m *handlerMockUserRepo) FindByID(id uint64) (*model.User, error)           { return nil, nil }
func (m *handlerMockUserRepo) PatchUpdate(id uint64, updates map[string]interface{}) error {
	return nil
}
func (m *handlerMockUserRepo) Delete(id uint64) error { return nil }

func handlerTestConfig() *config.Config {
	return &config.Config{
		JWTSecret:        "test-secret-key-for-handler-tests",
		JWTAccessExpiry:  15 * 60 * 1e9,
		JWTRefreshExpiry: 168 * 60 * 60 * 1e9,
	}
}

func mustHashPassword(pw string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func makePostRequest(url string, body interface{}) *http.Request {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, url, io.NopCloser(strings.NewReader(string(b))))
	req.Header.Set("Content-Type", "application/json")
	req.ContentLength = int64(len(b))
	return req
}

func makeGetRequest(url string) *http.Request {
	return httptest.NewRequest(http.MethodGet, url, nil)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	hashed := mustHashPassword("admin123")
	mockRepo := &handlerMockUserRepo{
		findByUsername: func(username string) (*model.User, error) {
			return &model.User{
				Status: 1,
			ID:           1,
				Username:     "admin",
				PasswordHash: hashed,
				Role:         "admin",
			}, nil
		},
	}

	authSvc := service.NewAuthService(mockRepo, handlerTestConfig())
	h := &Handler{svc: &service.Service{Auth: authSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makePostRequest("/api/v1/auth/login", dto.LoginRequest{
		Username: "admin",
		Password: "admin123",
	})

	h.Login(c)

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

func TestAuthHandler_Login_Failure_InvalidCredentials(t *testing.T) {
	mockRepo := &handlerMockUserRepo{
		findByUsername: func(username string) (*model.User, error) {
			return nil, errors.New("not found")
		},
	}

	authSvc := service.NewAuthService(mockRepo, handlerTestConfig())
	h := &Handler{svc: &service.Service{Auth: authSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makePostRequest("/api/v1/auth/login", dto.LoginRequest{
		Username: "baduser",
		Password: "badpass",
	})

	h.Login(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_Login_BadRequest(t *testing.T) {
	h := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	h.Login(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid body, got %d", w.Code)
	}
}

func TestAuthHandler_RefreshToken_NoCookie(t *testing.T) {
	h := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)

	h.RefreshToken(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401 when no cookie, got %d", w.Code)
	}
}

func TestAuthHandler_Login_EmptyJSON(t *testing.T) {
	h := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makePostRequest("/api/v1/auth/login", map[string]string{})

	h.Login(c)

	// Empty JSON {} misses required fields, so ShouldBindJSON fails with 400
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for empty JSON, got %d", w.Code)
	}
}

func TestAuthHandler_RefreshToken_InvalidCookie(t *testing.T) {
	authSvc := service.NewAuthService(&handlerMockUserRepo{}, handlerTestConfig())
	h := &Handler{svc: &service.Service{Auth: authSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/auth/refresh")
	c.Request.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: "some.invalid.token",
	})

	h.RefreshToken(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401 for invalid token, got %d", w.Code)
	}
}

func TestAuthHandler_RefreshToken_Success(t *testing.T) {
	hashed := mustHashPassword("admin123")
	mockRepo := &handlerMockUserRepo{
		findByUsername: func(username string) (*model.User, error) {
			return &model.User{
				Status: 1,
			ID:           1,
				Username:     "admin",
				PasswordHash: hashed,
				Role:         "admin",
			}, nil
		},
	}

	cfg := handlerTestConfig()
	authSvc := service.NewAuthService(mockRepo, cfg)

	// First, get a valid refresh token via login service call
	pair, err := authSvc.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("setup: login failed: %v", err)
	}

	h := &Handler{svc: &service.Service{Auth: authSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/auth/refresh")

	// Set the refresh token cookie as the handler expects
	c.Request.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: pair.RefreshToken,
	})

	h.RefreshToken(c)

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
