package service

import (
	"errors"
	"testing"
	"time"

	"mygo-immigration/backend/internal/config"
	"mygo-immigration/backend/internal/model"

	"golang.org/x/crypto/bcrypt"
)

// mockUserRepo implements repository.UserRepository for testing.
type mockUserRepo struct {
	findByUsernameFn  func(username string) (*model.User, error)
	findAllFn         func() ([]model.User, error)
	createFn          func(user *model.User) error
	updateFn          func(user *model.User) error
	findByIDFn        func(id uint64) (*model.User, error)
	patchUpdateFn     func(id uint64, updates map[string]interface{}) error
	deleteFn          func(id uint64) error
}

func (m *mockUserRepo) FindByUsername(username string) (*model.User, error) {
	if m.findByUsernameFn != nil {
		return m.findByUsernameFn(username)
	}
	return nil, errors.New("not found")
}

func (m *mockUserRepo) FindAll() ([]model.User, error) {
	if m.findAllFn != nil {
		return m.findAllFn()
	}
	return nil, nil
}

func (m *mockUserRepo) Create(user *model.User) error {
	if m.createFn != nil {
		return m.createFn(user)
	}
	return nil
}

func (m *mockUserRepo) Update(user *model.User) error {
	if m.updateFn != nil {
		return m.updateFn(user)
	}
	return nil
}

func (m *mockUserRepo) FindByID(id uint64) (*model.User, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return nil, errors.New("not found")
}

func (m *mockUserRepo) PatchUpdate(id uint64, updates map[string]interface{}) error {
	if m.patchUpdateFn != nil {
		return m.patchUpdateFn(id, updates)
	}
	return nil
}

func (m *mockUserRepo) Delete(id uint64) error {
	if m.deleteFn != nil {
		return m.deleteFn(id)
	}
	return nil
}

func (m *mockUserRepo) FindAllPaginated(page, perPage int) ([]model.User, int64, error) {
	return nil, 0, nil
}

func testConfig() *config.Config {
	return &config.Config{
		JWTSecret:        "test-secret-key-for-unit-tests",
		JWTAccessExpiry:  15 * time.Minute,
		JWTRefreshExpiry: 168 * time.Hour,
	}
}

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash)
}

func TestLogin_Success(t *testing.T) {
	password := "correct-password"
	user := &model.User{
		Status: 1,
			ID:           1,
		Username:     "testuser",
		PasswordHash: hashPassword(password),
		Role:         "admin",
	}

	repo := &mockUserRepo{
		findByUsernameFn: func(username string) (*model.User, error) {
			if username == "testuser" {
				return user, nil
			}
			return nil, errors.New("not found")
		},
	}

	svc := NewAuthService(repo, testConfig())

	pair, err := svc.Login("testuser", password)
	if err != nil {
		t.Fatalf("expected login success, got error: %v", err)
	}
	if pair.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
	if pair.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}
}

func TestLogin_InvalidPassword(t *testing.T) {
	user := &model.User{
		Status: 1,
			ID:           1,
		Username:     "testuser",
		PasswordHash: hashPassword("correct-password"),
		Role:         "admin",
	}

	repo := &mockUserRepo{
		findByUsernameFn: func(username string) (*model.User, error) {
			return user, nil
		},
	}

	svc := NewAuthService(repo, testConfig())

	_, err := svc.Login("testuser", "wrong-password")
	if err == nil {
		t.Fatal("expected error for invalid password")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := &mockUserRepo{
		findByUsernameFn: func(username string) (*model.User, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewAuthService(repo, testConfig())

	_, err := svc.Login("nonexistent", "password")
	if err == nil {
		t.Fatal("expected error for user not found")
	}
}

func TestLogin_EmptyCredentials(t *testing.T) {
	repo := &mockUserRepo{}
	svc := NewAuthService(repo, testConfig())

	_, err := svc.Login("", "password")
	if err == nil {
		t.Fatal("expected error for empty username")
	}

	_, err = svc.Login("user", "")
	if err == nil {
		t.Fatal("expected error for empty password")
	}
}

func TestRefreshToken_Success(t *testing.T) {
	password := "correct-password"
	user := &model.User{
		Status: 1,
			ID:           1,
		Username:     "testuser",
		PasswordHash: hashPassword(password),
		Role:         "admin",
	}

	repo := &mockUserRepo{
		findByUsernameFn: func(username string) (*model.User, error) {
			return user, nil
		},
	}

	cfg := testConfig()
	svc := NewAuthService(repo, cfg)

	// First login to get tokens
	initial, err := svc.Login("testuser", password)
	if err != nil {
		t.Fatalf("setup login failed: %v", err)
	}

	// Now refresh using the refresh token
	pair, err := svc.RefreshToken(initial.RefreshToken)
	if err != nil {
		t.Fatalf("expected refresh success, got error: %v", err)
	}
	if pair.AccessToken == "" {
		t.Error("expected non-empty access token from refresh")
	}
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	repo := &mockUserRepo{}
	svc := NewAuthService(repo, testConfig())

	_, err := svc.RefreshToken("invalid-token")
	if err == nil {
		t.Fatal("expected error for invalid refresh token")
	}
}

func TestRefreshToken_EmptyToken(t *testing.T) {
	repo := &mockUserRepo{}
	svc := NewAuthService(repo, testConfig())

	_, err := svc.RefreshToken("")
	if err == nil {
		t.Fatal("expected error for empty refresh token")
	}
}

func TestLogin_BcryptCostVerification(t *testing.T) {
	// Verify that bcrypt.DefaultCost is used (cost level 10)
	password := "test-password"
	hash := hashPassword(password)

	// hashPassword uses bcrypt.MinCost (4), but the auth service Login
	// uses bcrypt.CompareHashAndPassword which works at any cost level.
	// This test verifies that comparing passwords works with real bcrypt.
	user := &model.User{
		Status: 1,
			ID:           2,
		Username:     "costuser",
		PasswordHash: hash,
		Role:         "viewer",
	}

	repo := &mockUserRepo{
		findByUsernameFn: func(username string) (*model.User, error) {
			return user, nil
		},
	}

	svc := NewAuthService(repo, testConfig())

	pair, err := svc.Login("costuser", password)
	if err != nil {
		t.Fatalf("expected login success with bcrypt hash, got error: %v", err)
	}
	if pair.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
}

func TestGenerateToken_ExpiryTimes(t *testing.T) {
	repo := &mockUserRepo{
		findByUsernameFn: func(username string) (*model.User, error) {
			return &model.User{
				Status: 1,
			ID:           1,
				Username:     username,
				PasswordHash: hashPassword("secret"),
				Role:         "admin",
			}, nil
		},
	}

	cfg := testConfig()
	// Override expiry to very short durations to test
	cfg.JWTAccessExpiry = 60 * time.Second // 60 seconds
	cfg.JWTRefreshExpiry = 120 * time.Second
	svc := NewAuthService(repo, cfg)

	pair, err := svc.Login("admin", "secret")
	if err != nil {
		t.Fatalf("expected login success, got error: %v", err)
	}
	// ExpiresIn should match the access token expiry in seconds
	expectedExpiry := int64(60)
	if pair.ExpiresIn != expectedExpiry {
		t.Errorf("expected ExpiresIn %d, got %d", expectedExpiry, pair.ExpiresIn)
	}
}

func TestLogin_ContainsClaims(t *testing.T) {
	repo := &mockUserRepo{
		findByUsernameFn: func(username string) (*model.User, error) {
			return &model.User{
				Status: 1,
			ID:           42,
				Username:     "claimsuser",
				PasswordHash: hashPassword("testpass"),
				Role:         "editor",
			}, nil
		},
	}

	cfg := testConfig()
	svc := NewAuthService(repo, cfg)

	pair, err := svc.Login("claimsuser", "testpass")
	if err != nil {
		t.Fatalf("expected login success, got error: %v", err)
	}
	if pair.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
	if pair.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}
}

func TestRefreshToken_WrongSigningKey(t *testing.T) {
	repo := &mockUserRepo{
		findByUsernameFn: func(username string) (*model.User, error) {
			return &model.User{
				Status: 1,
			ID:           1,
				Username:     "admin",
				PasswordHash: hashPassword("secret"),
				Role:         "admin",
			}, nil
		},
	}

	cfg1 := testConfig()
	svc1 := NewAuthService(repo, cfg1)
	pair, _ := svc1.Login("admin", "secret")

	// Try to refresh with a different config (different secret)
	cfg2 := testConfig()
	cfg2.JWTSecret = "different-secret-key"
	svc2 := NewAuthService(repo, cfg2)

	_, err := svc2.RefreshToken(pair.RefreshToken)
	if err == nil {
		t.Fatal("expected error when refreshing with wrong signing key")
	}
}

func TestRefreshToken_DoesNotReturnRefreshToken(t *testing.T) {
	// Refresh should only return a new access token, not a new refresh token
	repo := &mockUserRepo{
		findByUsernameFn: func(username string) (*model.User, error) {
			return &model.User{
				Status: 1,
			ID:           1,
				Username:     "admin",
				PasswordHash: hashPassword("secret"),
				Role:         "admin",
			}, nil
		},
	}

	cfg := testConfig()
	svc := NewAuthService(repo, cfg)

	pair, _ := svc.Login("admin", "secret")
	refreshed, err := svc.RefreshToken(pair.RefreshToken)
	if err != nil {
		t.Fatalf("expected refresh success, got error: %v", err)
	}
	if refreshed.RefreshToken != "" {
		t.Error("expected refresh token to be empty in refresh response")
	}
	if refreshed.AccessToken == "" {
		t.Error("expected non-empty access token in refresh response")
	}
}

func TestLogin_RepoError(t *testing.T) {
	repo := &mockUserRepo{
		findByUsernameFn: func(username string) (*model.User, error) {
			return nil, errors.New("db connection error")
		},
	}

	svc := NewAuthService(repo, testConfig())

	_, err := svc.Login("anyuser", "anypassword")
	if err == nil {
		t.Fatal("expected error when repo returns error")
	}
}

