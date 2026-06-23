package service

import (
	"errors"
	"fmt"
	"regexp"
	"sort"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

const protectedAdminUsername = "admin"

var roleCodePattern = regexp.MustCompile(`^[a-z][a-z0-9_-]{1,63}$`)

type RBACService struct {
	repo     repository.RBACRepository
	userRepo repository.UserRepository
}

func NewRBACService(repo repository.RBACRepository, userRepo repository.UserRepository) *RBACService {
	return &RBACService{repo: repo, userRepo: userRepo}
}

func (s *RBACService) ListPermissions() ([]model.Permission, error) {
	permissions, err := s.repo.FindPermissions()
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}
	return permissions, nil
}

func (s *RBACService) ListRoles() ([]model.Role, error) {
	roles, err := s.repo.FindRoles(repository.RoleFilter{})
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	return roles, nil
}

type RoleDetail struct {
	model.Role
	PermissionCodes []string `json:"permission_codes"`
}

func (s *RBACService) GetRole(id uint64) (*RoleDetail, error) {
	role, err := s.repo.FindRoleByID(id)
	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}
	permissionCodes, err := s.repo.FindRolePermissionCodes(id)
	if err != nil {
		return nil, fmt.Errorf("failed to load role permissions: %w", err)
	}
	return &RoleDetail{Role: *role, PermissionCodes: permissionCodes}, nil
}

func (s *RBACService) CreateRole(req dto.CreateRoleRequest) (*RoleDetail, error) {
	if !roleCodePattern.MatchString(req.Code) {
		return nil, errors.New("invalid role code")
	}
	if req.Name == "" {
		return nil, errors.New("role name is required")
	}

	status := int8(1)
	if req.Status != nil {
		status = *req.Status
	}

	role := &model.Role{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Status:      status,
		IsSystem:    false,
	}
	if err := s.repo.CreateRole(role); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}
	if err := s.ReplaceRolePermissions(role.ID, req.PermissionIDs); err != nil {
		return nil, err
	}
	return s.GetRole(role.ID)
}

func (s *RBACService) UpdateRole(id uint64, req dto.UpdateRoleRequest) (*RoleDetail, error) {
	role, err := s.repo.FindRoleByID(id)
	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}
	if role.Code == protectedAdminUsername && req.PermissionIDs != nil {
		return nil, errors.New("admin role permissions cannot be changed")
	}
	if req.Name != "" {
		role.Name = req.Name
	}
	role.Description = req.Description
	if req.Status != nil {
		role.Status = *req.Status
	}
	if err := s.repo.UpdateRole(role); err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}
	if req.PermissionIDs != nil {
		if err := s.ReplaceRolePermissions(id, req.PermissionIDs); err != nil {
			return nil, err
		}
	}
	return s.GetRole(id)
}

func (s *RBACService) DeleteRole(id uint64) error {
	role, err := s.repo.FindRoleByID(id)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}
	if role.IsSystem {
		return errors.New("system role cannot be deleted")
	}
	if err := s.repo.DeleteRole(id); err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}
	return nil
}

func (s *RBACService) ReplaceRolePermissions(roleID uint64, permissionCodes []string) error {
	role, err := s.repo.FindRoleByID(roleID)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}
	if role.Code == protectedAdminUsername {
		return errors.New("admin role permissions cannot be changed")
	}
	permissionCodes = uniqueStrings(permissionCodes)
	if err := s.ensurePermissionCodesExist(permissionCodes); err != nil {
		return err
	}
	if err := s.repo.ReplaceRolePermissions(roleID, permissionCodes); err != nil {
		return fmt.Errorf("failed to save role permissions: %w", err)
	}
	return nil
}

func (s *RBACService) EffectivePermissions(userID uint64) ([]string, error) {
	if userID == 0 {
		return nil, errors.New("user id is required")
	}
	codes, err := s.repo.FindEffectivePermissionCodes(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to load permissions: %w", err)
	}
	return codes, nil
}

func (s *RBACService) HasPermission(userID uint64, permission string) (bool, error) {
	codes, err := s.EffectivePermissions(userID)
	if err != nil {
		return false, err
	}
	for _, code := range codes {
		if code == permission {
			return true, nil
		}
	}
	return false, nil
}

func (s *RBACService) UserOverrides(userID uint64) ([]model.UserPermissionOverride, error) {
	overrides, err := s.repo.FindUserPermissionOverrides(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to load user permission overrides: %w", err)
	}
	return overrides, nil
}

func (s *RBACService) ReplaceUserOverrides(userID uint64, overrides []dto.PermissionOverrideRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	if user.Username == protectedAdminUsername {
		for _, override := range overrides {
			if override.Effect == "deny" && (override.PermissionCode == "users:write" || override.PermissionCode == "roles:write") {
				return errors.New("protected admin must keep user and role management permissions")
			}
		}
	}

	inputs := make([]repository.PermissionOverrideInput, 0, len(overrides))
	codes := make([]string, 0, len(overrides))
	for _, override := range overrides {
		if override.PermissionCode == "" {
			continue
		}
		if override.Effect != "allow" && override.Effect != "deny" {
			return errors.New("invalid permission override effect")
		}
		codes = append(codes, override.PermissionCode)
		inputs = append(inputs, repository.PermissionOverrideInput{
			PermissionCode: override.PermissionCode,
			Effect:         override.Effect,
		})
	}
	if err := s.ensurePermissionCodesExist(codes); err != nil {
		return err
	}
	if err := s.repo.ReplaceUserPermissionOverrides(userID, inputs); err != nil {
		return fmt.Errorf("failed to save user permission overrides: %w", err)
	}
	return nil
}

func (s *RBACService) ensurePermissionCodesExist(codes []string) error {
	codes = uniqueStrings(codes)
	if len(codes) == 0 {
		return nil
	}
	permissions, err := s.repo.FindPermissionsByCodes(codes)
	if err != nil {
		return fmt.Errorf("failed to validate permissions: %w", err)
	}
	if len(permissions) != len(codes) {
		return errors.New("unknown permission code")
	}
	return nil
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]bool, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}
