package service

import (
	"errors"
	"fmt"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

// UserService handles business logic for user management.
type UserService struct {
	repo     repository.UserRepository
	rbacRepo repository.RBACRepository
}

// NewUserService creates a new UserService with the given repository.
func NewUserService(repo repository.UserRepository, rbacRepo ...repository.RBACRepository) *UserService {
	var rr repository.RBACRepository
	if len(rbacRepo) > 0 {
		rr = rbacRepo[0]
	}
	return &UserService{repo: repo, rbacRepo: rr}
}

// List returns users with optional filters and pagination.
func (s *UserService) List(req dto.UserListRequest) ([]model.User, int64, error) {
	users, total, err := s.repo.FindAll(repository.UserFilter{
		Role:     req.Role,
		Status:   req.Status,
		Username: req.Username,
		Page:     req.Page,
		PerPage:  req.PerPage,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	return users, total, nil
}

// Create creates a new user with a bcrypt-hashed password.
func (s *UserService) Create(username, password, displayName, role string, overrides ...[]dto.PermissionOverrideRequest) (*model.User, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}
	if role == "" {
		role = "viewer"
	}
	if s.rbacRepo != nil {
		if _, err := s.rbacRepo.FindRoleByCode(role); err != nil {
			return nil, fmt.Errorf("role not found: %w", err)
		}
	}
	var overrideInputs []repository.PermissionOverrideInput
	if s.rbacRepo != nil && len(overrides) > 0 {
		var err error
		overrideInputs, err = buildOverrideInputs(overrides[0])
		if err != nil {
			return nil, err
		}
		if err := validateOverridePermissions(s.rbacRepo, overrideInputs); err != nil {
			return nil, err
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		Username:     username,
		PasswordHash: string(hash),
		DisplayName:  displayName,
		Role:         role,
		Status:       1,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	if s.rbacRepo != nil && len(overrides) > 0 {
		if err := s.rbacRepo.ReplaceUserPermissionOverrides(user.ID, overrideInputs); err != nil {
			return nil, fmt.Errorf("failed to save user permission overrides: %w", err)
		}
	}
	return user, nil
}

func buildOverrideInputs(overrides []dto.PermissionOverrideRequest) ([]repository.PermissionOverrideInput, error) {
	inputs := make([]repository.PermissionOverrideInput, 0, len(overrides))
	for _, override := range overrides {
		if override.PermissionCode == "" {
			continue
		}
		if override.Effect != "allow" && override.Effect != "deny" {
			return nil, errors.New("invalid permission override effect")
		}
		inputs = append(inputs, repository.PermissionOverrideInput{
			PermissionCode: override.PermissionCode,
			Effect:         override.Effect,
		})
	}
	return inputs, nil
}

func validateOverridePermissions(repo repository.RBACRepository, inputs []repository.PermissionOverrideInput) error {
	codes := make([]string, 0, len(inputs))
	for _, input := range inputs {
		codes = append(codes, input.PermissionCode)
	}
	permissions, err := repo.FindPermissionsByCodes(codes)
	if err != nil {
		return fmt.Errorf("failed to validate permissions: %w", err)
	}
	if len(permissions) != len(codes) {
		return errors.New("unknown permission code")
	}
	return nil
}

// Update applies partial updates to a user. If password is set, it is bcrypt-hashed.
func (s *UserService) Update(id uint64, req dto.UpdateUserRequest) (*model.User, error) {
	if id == 0 {
		return nil, errors.New("user id is required")
	}

	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if req.DisplayName != "" {
		existing.DisplayName = req.DisplayName
	}
	if req.Role != "" {
		if existing.Username == protectedAdminUsername && req.Role != protectedAdminUsername {
			return nil, errors.New("protected admin role cannot be changed")
		}
		if s.rbacRepo != nil {
			if _, err := s.rbacRepo.FindRoleByCode(req.Role); err != nil {
				return nil, fmt.Errorf("role not found: %w", err)
			}
		}
		existing.Role = req.Role
	}
	if req.Status != nil {
		if existing.Username == protectedAdminUsername && *req.Status != 1 {
			return nil, errors.New("protected admin cannot be disabled")
		}
		existing.Status = *req.Status
	}
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		existing.PasswordHash = string(hash)
	}

	if err := s.repo.Update(existing); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	if s.rbacRepo != nil && req.PermissionOverrides != nil {
		for _, override := range req.PermissionOverrides {
			if existing.Username == protectedAdminUsername && override.Effect == "deny" && (override.PermissionCode == "users:write" || override.PermissionCode == "roles:write") {
				return nil, errors.New("protected admin must keep user and role management permissions")
			}
		}
		inputs, err := buildOverrideInputs(req.PermissionOverrides)
		if err != nil {
			return nil, err
		}
		if err := validateOverridePermissions(s.rbacRepo, inputs); err != nil {
			return nil, err
		}
		if err := s.rbacRepo.ReplaceUserPermissionOverrides(id, inputs); err != nil {
			return nil, fmt.Errorf("failed to save user permission overrides: %w", err)
		}
	}

	return existing, nil
}

// ChangePassword updates the current user's password after verifying the old password.
func (s *UserService) ChangePassword(userID uint64, oldPassword, newPassword string) error {
	if userID == 0 {
		return errors.New("user id is required")
	}
	if oldPassword == "" {
		return errors.New("old password is required")
	}
	if newPassword == "" {
		return errors.New("new password is required")
	}

	user, err := s.repo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("old password is incorrect")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.PasswordHash = string(hash)

	if err := s.repo.Update(user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

// Delete hard-deletes a user by ID.
func (s *UserService) Delete(id uint64) error {
	if id == 0 {
		return errors.New("user id is required")
	}
	if s.rbacRepo != nil {
		existing, err := s.repo.FindByID(id)
		if err != nil {
			return fmt.Errorf("user not found: %w", err)
		}
		if existing.Username == protectedAdminUsername {
			return errors.New("protected admin cannot be deleted")
		}
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// GetByID returns a user by ID.
func (s *UserService) GetByID(id uint64) (*model.User, error) {
	if id == 0 {
		return nil, errors.New("user id is required")
	}
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

// FindByUsername returns a user by username (used by auth).
func (s *UserService) FindByUsername(username string) (*model.User, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}
