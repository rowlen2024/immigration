package repository

import (
	"sort"

	"mygo-immigration/backend/internal/model"

	"gorm.io/gorm"
)

type RBACRepo struct {
	db *gorm.DB
}

func (r *RBACRepo) FindPermissions() ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.Order("sort_order asc, id asc").Find(&permissions).Error
	return permissions, err
}

func (r *RBACRepo) FindPermissionsByCodes(codes []string) ([]model.Permission, error) {
	var permissions []model.Permission
	if len(codes) == 0 {
		return permissions, nil
	}
	err := r.db.Where("code IN ?", codes).Find(&permissions).Error
	return permissions, err
}

func (r *RBACRepo) FindRoles(filter RoleFilter) ([]model.Role, error) {
	var roles []model.Role
	q := r.db.Model(&model.Role{})
	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
	}
	err := q.Order("is_system desc, id asc").Find(&roles).Error
	return roles, err
}

func (r *RBACRepo) FindRoleByID(id uint64) (*model.Role, error) {
	var role model.Role
	if err := r.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RBACRepo) FindRoleByCode(code string) (*model.Role, error) {
	var role model.Role
	if err := r.db.Where("code = ?", code).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RBACRepo) CreateRole(role *model.Role) error {
	return r.db.Create(role).Error
}

func (r *RBACRepo) UpdateRole(role *model.Role) error {
	return r.db.Omit("created_at").Save(role).Error
}

func (r *RBACRepo) DeleteRole(id uint64) error {
	return r.db.Unscoped().Delete(&model.Role{}, id).Error
}

func (r *RBACRepo) FindRolePermissionCodes(roleID uint64) ([]string, error) {
	var codes []string
	err := r.db.Table("permissions").
		Select("permissions.code").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Order("permissions.sort_order asc, permissions.id asc").
		Scan(&codes).Error
	return codes, err
}

func (r *RBACRepo) ReplaceRolePermissions(roleID uint64, permissionCodes []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
			return err
		}

		var permissions []model.Permission
		if len(permissionCodes) > 0 {
			if err := tx.Where("code IN ?", permissionCodes).Find(&permissions).Error; err != nil {
				return err
			}
		}

		items := make([]model.RolePermission, 0, len(permissions))
		for _, permission := range permissions {
			items = append(items, model.RolePermission{RoleID: roleID, PermissionID: permission.ID})
		}
		if len(items) == 0 {
			return nil
		}
		return tx.Create(&items).Error
	})
}

func (r *RBACRepo) FindUserPermissionOverrides(userID uint64) ([]model.UserPermissionOverride, error) {
	var overrides []model.UserPermissionOverride
	err := r.db.Preload("Permission").
		Where("user_id = ?", userID).
		Order("id asc").
		Find(&overrides).Error
	return overrides, err
}

func (r *RBACRepo) ReplaceUserPermissionOverrides(userID uint64, overrides []PermissionOverrideInput) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserPermissionOverride{}).Error; err != nil {
			return err
		}
		if len(overrides) == 0 {
			return nil
		}

		codes := make([]string, 0, len(overrides))
		effectsByCode := make(map[string]string, len(overrides))
		for _, override := range overrides {
			if override.PermissionCode == "" || override.Effect == "" {
				continue
			}
			codes = append(codes, override.PermissionCode)
			effectsByCode[override.PermissionCode] = override.Effect
		}
		if len(codes) == 0 {
			return nil
		}

		var permissions []model.Permission
		if err := tx.Where("code IN ?", codes).Find(&permissions).Error; err != nil {
			return err
		}

		items := make([]model.UserPermissionOverride, 0, len(permissions))
		for _, permission := range permissions {
			effect := effectsByCode[permission.Code]
			if effect != "allow" && effect != "deny" {
				continue
			}
			items = append(items, model.UserPermissionOverride{
				UserID:       userID,
				PermissionID: permission.ID,
				Effect:       effect,
			})
		}
		if len(items) == 0 {
			return nil
		}
		return tx.Create(&items).Error
	})
}

func (r *RBACRepo) FindEffectivePermissionCodes(userID uint64) ([]string, error) {
	var user model.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	var roleCodes []string
	if user.Status == 1 {
		err := r.db.Table("permissions").
			Select("permissions.code").
			Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
			Joins("JOIN roles ON roles.id = role_permissions.role_id").
			Where("roles.code = ? AND roles.status = 1", user.Role).
			Scan(&roleCodes).Error
		if err != nil {
			return nil, err
		}
	}

	allowed := make(map[string]bool, len(roleCodes))
	for _, code := range roleCodes {
		allowed[code] = true
	}

	var overrides []struct {
		Code   string
		Effect string
	}
	if err := r.db.Table("user_permission_overrides").
		Select("permissions.code, user_permission_overrides.effect").
		Joins("JOIN permissions ON permissions.id = user_permission_overrides.permission_id").
		Where("user_permission_overrides.user_id = ?", userID).
		Scan(&overrides).Error; err != nil {
		return nil, err
	}

	for _, override := range overrides {
		if override.Effect == "allow" {
			allowed[override.Code] = true
		}
		if override.Effect == "deny" {
			delete(allowed, override.Code)
		}
	}

	codes := make([]string, 0, len(allowed))
	for code := range allowed {
		codes = append(codes, code)
	}
	sort.Strings(codes)
	return codes, nil
}
