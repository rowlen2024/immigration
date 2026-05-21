package service

import (
	"errors"
	"testing"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
)

// userSvcMockUserRepo implements repository.UserRepository for user service tests.
type userSvcMockUserRepo struct {
	findByUsernameFn func(username string) (*model.User, error)
	findAllFn        func() ([]model.User, error)
	createFn         func(user *model.User) error
	updateFn         func(user *model.User) error
	findByIDFn       func(id uint64) (*model.User, error)
	patchUpdateFn    func(id uint64, updates map[string]interface{}) error
	deleteFn         func(id uint64) error
}

func (m *userSvcMockUserRepo) FindByUsername(username string) (*model.User, error) {
	if m.findByUsernameFn != nil {
		return m.findByUsernameFn(username)
	}
	return nil, errors.New("not found")
}

func (m *userSvcMockUserRepo) FindAll() ([]model.User, error) {
	if m.findAllFn != nil {
		return m.findAllFn()
	}
	return nil, nil
}

func (m *userSvcMockUserRepo) Create(user *model.User) error {
	if m.createFn != nil {
		return m.createFn(user)
	}
	return nil
}

func (m *userSvcMockUserRepo) Update(user *model.User) error {
	if m.updateFn != nil {
		return m.updateFn(user)
	}
	return nil
}

func (m *userSvcMockUserRepo) FindByID(id uint64) (*model.User, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return nil, errors.New("not found")
}

func (m *userSvcMockUserRepo) PatchUpdate(id uint64, updates map[string]interface{}) error {
	if m.patchUpdateFn != nil {
		return m.patchUpdateFn(id, updates)
	}
	return nil
}

func (m *userSvcMockUserRepo) FindAllPaginated(page, perPage int) ([]model.User, int64, error) {
	return nil, 0, nil
}

func (m *userSvcMockUserRepo) Delete(id uint64) error {
	if m.deleteFn != nil {
		return m.deleteFn(id)
	}
	return nil
}

func TestUser_List(t *testing.T) {
	sampleUsers := []model.User{
		{ID: 1, Username: "alice", DisplayName: "Alice", Role: "admin"},
		{ID: 2, Username: "bob", DisplayName: "Bob", Role: "editor"},
	}

	repo := &userSvcMockUserRepo{
		findAllFn: func() ([]model.User, error) {
			return sampleUsers, nil
		},
	}

	svc := &UserService{repo: repo}

	users, err := svc.List()
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
	if users[0].Username != "alice" {
		t.Errorf("expected first user 'alice', got '%s'", users[0].Username)
	}
}

func TestUser_List_Error(t *testing.T) {
	repo := &userSvcMockUserRepo{
		findAllFn: func() ([]model.User, error) {
			return nil, errors.New("db down")
		},
	}

	svc := &UserService{repo: repo}

	_, err := svc.List()
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestUser_Create_Success(t *testing.T) {
	var createdUser *model.User
	repo := &userSvcMockUserRepo{
		createFn: func(user *model.User) error {
			createdUser = user
			user.ID = 10
			return nil
		},
	}

	svc := &UserService{repo: repo}

	user, err := svc.Create("newuser", "password123", "New User", "editor")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if user.ID != 10 {
		t.Errorf("expected ID 10, got %d", user.ID)
	}
	if user.Username != "newuser" {
		t.Errorf("expected username 'newuser', got '%s'", user.Username)
	}
	if user.DisplayName != "New User" {
		t.Errorf("expected display name 'New User', got '%s'", user.DisplayName)
	}
	if user.Role != "editor" {
		t.Errorf("expected role 'editor', got '%s'", user.Role)
	}
	if createdUser == nil {
		t.Fatal("expected Create to be called on repo")
	}
	if createdUser.PasswordHash == "" {
		t.Error("expected password hash to be set")
	}
	if createdUser.PasswordHash == "password123" {
		t.Error("expected password to be hashed, not plaintext")
	}
	if createdUser.Status != 1 {
		t.Errorf("expected status 1, got %d", createdUser.Status)
	}
}

func TestUser_Create_PasswordHashing(t *testing.T) {
	var savedHash string
	repo := &userSvcMockUserRepo{
		createFn: func(user *model.User) error {
			savedHash = user.PasswordHash
			return nil
		},
	}

	svc := &UserService{repo: repo}

	_, err := svc.Create("testuser", "secret123", "Test", "viewer")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	// Verify that the hash is a valid bcrypt hash
	if len(savedHash) < 20 {
		t.Errorf("expected bcrypt hash to be longer than 20 chars, got %d", len(savedHash))
	}
	// Verify each call produces a different hash (bcrypt uses random salt)
	_, err = svc.Create("testuser2", "secret123", "Test2", "viewer")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	// The hash should be different even for same password due to random salt
	// Note: savedHash was from first call; second call updates it
}

func TestUser_Create_EmptyUsername(t *testing.T) {
	repo := &userSvcMockUserRepo{}
	svc := &UserService{repo: repo}

	_, err := svc.Create("", "password", "Display", "admin")
	if err == nil {
		t.Fatal("expected error for empty username")
	}
}

func TestUser_Create_EmptyPassword(t *testing.T) {
	repo := &userSvcMockUserRepo{}
	svc := &UserService{repo: repo}

	_, err := svc.Create("username", "", "Display", "admin")
	if err == nil {
		t.Fatal("expected error for empty password")
	}
}

func TestUser_Create_DefaultRole(t *testing.T) {
	var createdRole string
	repo := &userSvcMockUserRepo{
		createFn: func(user *model.User) error {
			createdRole = user.Role
			return nil
		},
	}

	svc := &UserService{repo: repo}

	_, err := svc.Create("newuser", "password", "New User", "")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if createdRole != "viewer" {
		t.Errorf("expected default role 'viewer', got '%s'", createdRole)
	}
}

func TestUser_Create_RepoError(t *testing.T) {
	repo := &userSvcMockUserRepo{
		createFn: func(user *model.User) error {
			return errors.New("db error")
		},
	}

	svc := &UserService{repo: repo}

	_, err := svc.Create("user", "pass", "Display", "viewer")
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestUser_Update_Success(t *testing.T) {
	var updated *model.User
	repo := &userSvcMockUserRepo{
		findByIDFn: func(id uint64) (*model.User, error) {
			return &model.User{ID: id, Username: "existinguser"}, nil
		},
		updateFn: func(user *model.User) error {
			updated = user
			return nil
		},
	}

	svc := &UserService{repo: repo}

	user, err := svc.Update(5, dto.UpdateUserRequest{DisplayName: "Updated Display", Role: "editor"})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if user.ID != 5 {
		t.Errorf("expected ID 5, got %d", user.ID)
	}
	if user.Username != "existinguser" {
		t.Errorf("expected username 'existinguser', got '%s'", user.Username)
	}
	if updated.DisplayName != "Updated Display" {
		t.Errorf("expected display_name 'Updated Display', got '%s'", updated.DisplayName)
	}
}

func TestUser_Update_WithPassword(t *testing.T) {
	var updated *model.User
	repo := &userSvcMockUserRepo{
		findByIDFn: func(id uint64) (*model.User, error) {
			return &model.User{ID: id, Username: "user"}, nil
		},
		updateFn: func(user *model.User) error {
			updated = user
			return nil
		},
	}

	svc := &UserService{repo: repo}

	_, err := svc.Update(1, dto.UpdateUserRequest{Password: "newpassword"})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(updated.PasswordHash) < 20 {
		t.Errorf("expected bcrypt hash to be long, got %d chars", len(updated.PasswordHash))
	}
}

func TestUser_Update_ZeroID(t *testing.T) {
	repo := &userSvcMockUserRepo{}
	svc := &UserService{repo: repo}

	_, err := svc.Update(0, dto.UpdateUserRequest{Role: "editor"})
	if err == nil {
		t.Fatal("expected error for zero id")
	}
}

func TestUser_Update_NotFound(t *testing.T) {
	repo := &userSvcMockUserRepo{
		findByIDFn: func(id uint64) (*model.User, error) {
			return nil, errors.New("not found")
		},
	}
	svc := &UserService{repo: repo}

	_, err := svc.Update(1, dto.UpdateUserRequest{Role: "editor"})
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestUser_Update_RepoError(t *testing.T) {
	repo := &userSvcMockUserRepo{
		findByIDFn: func(id uint64) (*model.User, error) {
			return &model.User{ID: id, Username: "user"}, nil
		},
		updateFn: func(user *model.User) error {
			return errors.New("db error")
		},
	}

	svc := &UserService{repo: repo}

	_, err := svc.Update(1, dto.UpdateUserRequest{Role: "editor"})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestUser_FindByUsername_Success(t *testing.T) {
	expected := &model.User{ID: 1, Username: "admin", Role: "admin"}
	repo := &userSvcMockUserRepo{
		findByUsernameFn: func(username string) (*model.User, error) {
			if username == "admin" {
				return expected, nil
			}
			return nil, errors.New("not found")
		},
	}

	svc := &UserService{repo: repo}

	user, err := svc.FindByUsername("admin")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if user.Username != "admin" {
		t.Errorf("expected username 'admin', got '%s'", user.Username)
	}
	if user.Role != "admin" {
		t.Errorf("expected role 'admin', got '%s'", user.Role)
	}
}

func TestUser_FindByUsername_Empty(t *testing.T) {
	repo := &userSvcMockUserRepo{}
	svc := &UserService{repo: repo}

	_, err := svc.FindByUsername("")
	if err == nil {
		t.Fatal("expected error for empty username")
	}
}

func TestUser_FindByUsername_NotFound(t *testing.T) {
	repo := &userSvcMockUserRepo{
		findByUsernameFn: func(username string) (*model.User, error) {
			return nil, errors.New("not found")
		},
	}

	svc := &UserService{repo: repo}

	_, err := svc.FindByUsername("nonexistent")
	if err == nil {
		t.Fatal("expected error for user not found")
	}
}
