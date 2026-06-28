package dto

// UserListRequest 用户列表查询参数
type UserListRequest struct {
	PaginationRequest
	Role     string `form:"role"`
	Status   *int8  `form:"status"`
	Username string `form:"username"`
}

// CreateUserRequest 用户创建请求
type CreateUserRequest struct {
	Username            string                      `json:"username" binding:"required"`
	Password            string                      `json:"password" binding:"required"`
	DisplayName         string                      `json:"display_name"`
	Role                string                      `json:"role"`
	PermissionOverrides []PermissionOverrideRequest `json:"permission_overrides"`
}

// ChangePasswordRequest holds current-user password change payload.
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
