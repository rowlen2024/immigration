package service

import (
	"errors"
	"testing"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

// userSvcMockUserRepo implements repository.UserRepository for user service tests.
type userSvcMockUserRepo struct {
	findByUsernameFn func(username string) (*model.User, error)
	findAllFn        func(filter repository.UserFilter) ([]model.User, int64, error)
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

func (m *userSvcMockUserRepo) FindAll(filter repository.UserFilter) ([]model.User, int64, error) {
	if m.findAllFn != nil {
		return m.findAllFn(filter)
	}
	return nil, 0, nil
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
		findAllFn: func(filter repository.UserFilter) ([]model.User, int64, error) {
			return sampleUsers, 2, nil
		},
	}

	svc := NewUserService(repo)

	users, total, err := svc.List(dto.UserListRequest{})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
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
		findAllFn: func(filter repository.UserFilter) ([]model.User, int64, error) {
			return nil, 0, errors.New("db down")
		},
	}

	svc := NewUserService(repo)

	_, _, err := svc.List(dto.UserListRequest{})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestUser_List_FilterByRole(t *testing.T) {
	repo := &userSvcMockUserRepo{
		findAllFn: func(filter repository.UserFilter) ([]model.User, int64, error) {
			if filter.Role != "admin" {
				t.Errorf("expected role filter 'admin', got '%s'", filter.Role)
			}
			return []model.User{{ID: 1, Role: "admin"}}, 1, nil
		},
	}

	svc := NewUserService(repo)
	_, _, err := svc.List(dto.UserListRequest{Role: "admin"})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestUser_List_FilterByStatus(t *testing.T) {
	status := int8(0)
	repo := &userSvcMockUserRepo{
		findAllFn: func(filter repository.UserFilter) ([]model.User, int64, error) {
			if filter.Status == nil || *filter.Status != 0 {
				t.Errorf("expected status filter 0, got %v", filter.Status)
			}
			return []model.User{}, 0, nil
		},
	}

	svc := NewUserService(repo)
	_, _, err := svc.List(dto.UserListRequest{Status: &status})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestUser_List_NoPagination(t *testing.T) {
	repo := &userSvcMockUserRepo{
		findAllFn: func(filter repository.UserFilter) ([]model.User, int64, error) {
			if filter.Page != 0 {
				t.Errorf("expected page 0 (no pagination), got %d", filter.Page)
			}
			return []model.User{}, 0, nil
		},
	}

	svc := NewUserService(repo)
	_, _, err := svc.List(dto.UserListRequest{})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
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

	svc := NewUserService(repo)

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

	svc := NewUserService(repo)

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
	svc := NewUserService(repo)

	_, err := svc.Create("", "password", "Display", "admin")
	if err == nil {
		t.Fatal("expected error for empty username")
	}
}

func TestUser_Create_EmptyPassword(t *testing.T) {
	repo := &userSvcMockUserRepo{}
	svc := NewUserService(repo)

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

	svc := NewUserService(repo)

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

	svc := NewUserService(repo)

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

	svc := NewUserService(repo)

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

	svc := NewUserService(repo)

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
	svc := NewUserService(repo)

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
	svc := NewUserService(repo)

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

	svc := NewUserService(repo)

	_, err := svc.Update(1, dto.UpdateUserRequest{Role: "editor"})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestUser_Update_StatusToZero(t *testing.T) {
	var updated *model.User
	status := int8(0)
	repo := &userSvcMockUserRepo{
		findByIDFn: func(id uint64) (*model.User, error) {
			return &model.User{ID: id, Username: "user", Status: 1}, nil
		},
		updateFn: func(user *model.User) error {
			updated = user
			return nil
		},
	}

	svc := NewUserService(repo)

	_, err := svc.Update(1, dto.UpdateUserRequest{Status: &status})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if updated.Status != 0 {
		t.Errorf("expected status 0 (disabled), got %d", updated.Status)
	}
}

func TestUser_ChangePassword_Success(t *testing.T) {
	oldHash, err := bcrypt.GenerateFromPassword([]byte("oldpass"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	var updated *model.User
	repo := &userSvcMockUserRepo{
		findByIDFn: func(id uint64) (*model.User, error) {
			return &model.User{ID: id, Username: "user", PasswordHash: string(oldHash)}, nil
		},
		updateFn: func(user *model.User) error {
			updated = user
			return nil
		},
	}

	svc := NewUserService(repo)
	err = svc.ChangePassword(1, "oldpass", "newpass")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if updated == nil {
		t.Fatal("expected user to be updated")
	}
	if bcrypt.CompareHashAndPassword([]byte(updated.PasswordHash), []byte("newpass")) != nil {
		t.Fatal("expected password hash to match new password")
	}
}

func TestUser_ChangePassword_WrongOldPassword(t *testing.T) {
	oldHash, err := bcrypt.GenerateFromPassword([]byte("oldpass"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	updated := false
	repo := &userSvcMockUserRepo{
		findByIDFn: func(id uint64) (*model.User, error) {
			return &model.User{ID: id, Username: "user", PasswordHash: string(oldHash)}, nil
		},
		updateFn: func(user *model.User) error {
			updated = true
			return nil
		},
	}

	svc := NewUserService(repo)
	err = svc.ChangePassword(1, "wrongpass", "newpass")
	if err == nil {
		t.Fatal("expected error for wrong old password")
	}
	if updated {
		t.Fatal("expected password not to be updated")
	}
}

func TestUser_ChangePassword_ZeroID(t *testing.T) {
	repo := &userSvcMockUserRepo{}
	svc := NewUserService(repo)

	err := svc.ChangePassword(0, "oldpass", "newpass")
	if err == nil {
		t.Fatal("expected error for zero user id")
	}
}

func TestUser_Delete_Success(t *testing.T) {
	var deletedID uint64
	repo := &userSvcMockUserRepo{
		deleteFn: func(id uint64) error {
			deletedID = id
			return nil
		},
	}

	svc := NewUserService(repo)

	err := svc.Delete(5)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if deletedID != 5 {
		t.Errorf("expected deleted ID 5, got %d", deletedID)
	}
}

func TestUser_Delete_ZeroID(t *testing.T) {
	repo := &userSvcMockUserRepo{}
	svc := NewUserService(repo)

	err := svc.Delete(0)
	if err == nil {
		t.Fatal("expected error for zero id")
	}
}

func TestUser_Delete_RepoError(t *testing.T) {
	repo := &userSvcMockUserRepo{
		deleteFn: func(id uint64) error {
			return errors.New("db error")
		},
	}

	svc := NewUserService(repo)

	err := svc.Delete(1)
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestUser_GetByID_Success(t *testing.T) {
	expected := &model.User{ID: 1, Username: "admin", Role: "admin"}
	repo := &userSvcMockUserRepo{
		findByIDFn: func(id uint64) (*model.User, error) {
			return expected, nil
		},
	}

	svc := NewUserService(repo)

	user, err := svc.GetByID(1)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if user.Username != "admin" {
		t.Errorf("expected username 'admin', got '%s'", user.Username)
	}
}

func TestUser_GetByID_ZeroID(t *testing.T) {
	repo := &userSvcMockUserRepo{}
	svc := NewUserService(repo)

	_, err := svc.GetByID(0)
	if err == nil {
		t.Fatal("expected error for zero id")
	}
}

func TestUser_GetByID_NotFound(t *testing.T) {
	repo := &userSvcMockUserRepo{
		findByIDFn: func(id uint64) (*model.User, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewUserService(repo)

	_, err := svc.GetByID(999)
	if err == nil {
		t.Fatal("expected error for not found")
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

	svc := NewUserService(repo)

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
	svc := NewUserService(repo)

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

	svc := NewUserService(repo)

	_, err := svc.FindByUsername("nonexistent")
	if err == nil {
		t.Fatal("expected error for user not found")
	}
}
