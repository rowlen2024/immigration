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
	repo repository.UserRepository
}

// NewUserService creates a new UserService with the given repository.
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
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
func (s *UserService) Create(username, password, displayName, role string) (*model.User, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}
	if role == "" {
		role = "viewer"
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
	return user, nil
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
		existing.Role = req.Role
	}
	if req.Status != nil {
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

	return existing, nil
}

// Delete hard-deletes a user by ID.
func (s *UserService) Delete(id uint64) error {
	if id == 0 {
		return errors.New("user id is required")
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
