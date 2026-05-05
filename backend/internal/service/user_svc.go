package service

import (
	"errors"
	"fmt"

	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

// UserService handles business logic for user management.
type UserService struct {
	repo repository.UserRepository
}

// List returns all users.
func (s *UserService) List() ([]model.User, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

// ListPaginated returns users with pagination.
func (s *UserService) ListPaginated(page, perPage int) ([]model.User, int64, error) {
	users, total, err := s.repo.FindAllPaginated(page, perPage)
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

// Update applies partial updates to a user. If the "password" key is present, it is bcrypt-hashed.
func (s *UserService) Update(id uint64, updates map[string]interface{}) (*model.User, error) {
	if id == 0 {
		return nil, errors.New("user id is required")
	}
	if len(updates) == 0 {
		return nil, errors.New("no updates provided")
	}

	if password, ok := updates["password"]; ok {
		if pwStr, ok := password.(string); ok && pwStr != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(pwStr), bcrypt.DefaultCost)
			if err != nil {
				return nil, fmt.Errorf("failed to hash password: %w", err)
			}
			updates["password_hash"] = string(hash)
		}
		delete(updates, "password")
	}

	if err := s.repo.PatchUpdate(id, updates); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return s.repo.FindByID(id)
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
