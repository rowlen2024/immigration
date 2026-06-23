package dto

type PermissionOverrideRequest struct {
	PermissionCode string `json:"permission_code"`
	Effect         string `json:"effect"`
}

type CreateRoleRequest struct {
	Code          string   `json:"code" binding:"required"`
	Name          string   `json:"name" binding:"required"`
	Description   string   `json:"description"`
	Status        *int8    `json:"status"`
	PermissionIDs []string `json:"permission_codes"`
}

type UpdateRoleRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Status        *int8    `json:"status"`
	PermissionIDs []string `json:"permission_codes"`
}

type SaveRolePermissionsRequest struct {
	PermissionIDs []string `json:"permission_codes"`
}
