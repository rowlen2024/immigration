package service

import (
	"errors"
	"testing"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

type rbacSvcMockRepo struct {
	role                *model.Role
	updateRoleCalled    bool
	userOverridesCalled bool
}

func (m *rbacSvcMockRepo) FindPermissions() ([]model.Permission, error) { return nil, nil }
func (m *rbacSvcMockRepo) FindPermissionsByCodes(codes []string) ([]model.Permission, error) {
	permissions := make([]model.Permission, 0, len(codes))
	for i, code := range codes {
		permissions = append(permissions, model.Permission{ID: uint64(i + 1), Code: code})
	}
	return permissions, nil
}
func (m *rbacSvcMockRepo) FindRoles(filter repository.RoleFilter) ([]model.Role, error) {
	return nil, nil
}
func (m *rbacSvcMockRepo) FindRoleByID(id uint64) (*model.Role, error) {
	if m.role != nil {
		return m.role, nil
	}
	return &model.Role{ID: id, Code: "editor"}, nil
}
func (m *rbacSvcMockRepo) FindRoleByCode(code string) (*model.Role, error) {
	return &model.Role{ID: 1, Code: code}, nil
}
func (m *rbacSvcMockRepo) CreateRole(role *model.Role) error { return nil }
func (m *rbacSvcMockRepo) UpdateRole(role *model.Role) error {
	m.updateRoleCalled = true
	return nil
}
func (m *rbacSvcMockRepo) DeleteRole(id uint64) error { return nil }
func (m *rbacSvcMockRepo) FindRolePermissionCodes(roleID uint64) ([]string, error) {
	return nil, nil
}
func (m *rbacSvcMockRepo) ReplaceRolePermissions(roleID uint64, permissionCodes []string) error {
	return nil
}
func (m *rbacSvcMockRepo) FindUserPermissionOverrides(userID uint64) ([]model.UserPermissionOverride, error) {
	return nil, nil
}
func (m *rbacSvcMockRepo) ReplaceUserPermissionOverrides(userID uint64, overrides []repository.PermissionOverrideInput) error {
	m.userOverridesCalled = true
	return nil
}
func (m *rbacSvcMockRepo) FindEffectivePermissionCodes(userID uint64) ([]string, error) {
	return []string{"cases:read", "projects:read"}, nil
}

type rbacSvcMockUserRepo struct {
	user *model.User
}

func (m *rbacSvcMockUserRepo) FindByUsername(username string) (*model.User, error) {
	return nil, errors.New("not implemented")
}
func (m *rbacSvcMockUserRepo) FindAll(filter repository.UserFilter) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (m *rbacSvcMockUserRepo) Create(user *model.User) error { return nil }
func (m *rbacSvcMockUserRepo) Update(user *model.User) error { return nil }
func (m *rbacSvcMockUserRepo) FindByID(id uint64) (*model.User, error) {
	if m.user != nil {
		return m.user, nil
	}
	return &model.User{ID: id, Username: "editor", Role: "editor"}, nil
}
func (m *rbacSvcMockUserRepo) PatchUpdate(id uint64, updates map[string]interface{}) error {
	return nil
}
func (m *rbacSvcMockUserRepo) Delete(id uint64) error { return nil }

func TestRBACService_HasPermission(t *testing.T) {
	svc := NewRBACService(&rbacSvcMockRepo{}, &rbacSvcMockUserRepo{})

	ok, err := svc.HasPermission(1, "cases:read")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !ok {
		t.Fatal("expected cases:read to be allowed")
	}

	ok, err = svc.HasPermission(1, "roles:write")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ok {
		t.Fatal("expected roles:write to be denied")
	}
}

func TestRBACService_AdminRolePermissionsCannotBeChanged(t *testing.T) {
	repo := &rbacSvcMockRepo{role: &model.Role{ID: 1, Code: "admin"}}
	svc := NewRBACService(repo, &rbacSvcMockUserRepo{})

	err := svc.ReplaceRolePermissions(1, []string{"users:write", "roles:write"})
	if err == nil {
		t.Fatal("expected error when admin role permissions change")
	}
}

func TestRBACService_UpdateAdminRoleRejectsPermissionChangeBeforeSave(t *testing.T) {
	repo := &rbacSvcMockRepo{role: &model.Role{ID: 1, Code: "admin", Name: "Admin"}}
	svc := NewRBACService(repo, &rbacSvcMockUserRepo{})

	_, err := svc.UpdateRole(1, dto.UpdateRoleRequest{
		Name:          "Changed",
		PermissionIDs: []string{"users:write", "roles:write"},
	})
	if err == nil {
		t.Fatal("expected error when admin role permissions change")
	}
	if repo.updateRoleCalled {
		t.Fatal("expected admin role not to be saved before permission rejection")
	}
}

func TestRBACService_ProtectedAdminCannotDenyManagementPermissions(t *testing.T) {
	repo := &rbacSvcMockRepo{}
	userRepo := &rbacSvcMockUserRepo{user: &model.User{ID: 1, Username: "admin", Role: "admin"}}
	svc := NewRBACService(repo, userRepo)

	err := svc.ReplaceUserOverrides(1, []dto.PermissionOverrideRequest{
		{PermissionCode: "roles:write", Effect: "deny"},
	})
	if err == nil {
		t.Fatal("expected protected admin deny override to fail")
	}
	if repo.userOverridesCalled {
		t.Fatal("expected overrides not to be saved")
	}
}

func TestRBACService_SystemRoleCannotBeDeleted(t *testing.T) {
	repo := &rbacSvcMockRepo{role: &model.Role{ID: 1, Code: "viewer", IsSystem: true}}
	svc := NewRBACService(repo, &rbacSvcMockUserRepo{})

	if err := svc.DeleteRole(1); err == nil {
		t.Fatal("expected system role delete to fail")
	}
}
