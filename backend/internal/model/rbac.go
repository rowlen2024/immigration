package model

import "time"

type Role struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        string    `gorm:"uniqueIndex;size:64;not null" json:"code"`
	Name        string    `gorm:"size:128;not null" json:"name"`
	Description string    `gorm:"size:255;not null;default:''" json:"description"`
	Status      int8      `gorm:"not null;default:1" json:"status"`
	IsSystem    bool      `gorm:"not null;default:false" json:"is_system"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Role) TableName() string { return "roles" }

type Permission struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Code      string    `gorm:"uniqueIndex;size:96;not null" json:"code"`
	Name      string    `gorm:"size:128;not null" json:"name"`
	Module    string    `gorm:"size:64;not null;index" json:"module"`
	Action    string    `gorm:"size:32;not null" json:"action"`
	SortOrder int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Permission) TableName() string { return "permissions" }

type RolePermission struct {
	RoleID       uint64    `gorm:"primaryKey;autoIncrement:false" json:"role_id"`
	PermissionID uint64    `gorm:"primaryKey;autoIncrement:false" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func (RolePermission) TableName() string { return "role_permissions" }

type UserPermissionOverride struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint64     `gorm:"not null;index:uk_user_permission,unique" json:"user_id"`
	PermissionID uint64     `gorm:"not null;index:uk_user_permission,unique" json:"permission_id"`
	Effect       string     `gorm:"type:enum('allow','deny');not null" json:"effect"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Permission   Permission `gorm:"foreignKey:PermissionID" json:"permission"`
}

func (UserPermissionOverride) TableName() string { return "user_permission_overrides" }
