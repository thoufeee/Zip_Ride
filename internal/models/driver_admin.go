package models

import "time"

type DriverAdmin struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex" json:"email"`
	Password     string    `json:"-"`
	PasswordHash string    `json:"-"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Name         string    `json:"name"` // Keep for backward compatibility
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	// Relationships
	AccountRoles []DriverAdminAccountRole `gorm:"foreignKey:AdminID" json:"account_roles,omitempty"`
}

type DriverAdminPermission struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Key         string    `gorm:"uniqueIndex" json:"key"` // e.g., driver.view, driver.approve
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DriverAdminRole struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex" json:"name"` // e.g., driver_admin, reviewer
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationships
	RolePermissions []DriverAdminRolePermission `gorm:"foreignKey:RoleID" json:"role_permissions,omitempty"`
}

type DriverAdminRolePermission struct {
	RoleID       uint      `gorm:"index" json:"role_id"`
	PermissionID uint      `gorm:"index" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
	
	// Relationships
	Permission DriverAdminPermission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
}

type DriverAdminAccountRole struct {
	AdminID     uint      `gorm:"index" json:"admin_id"`
	DriverAdminID uint    `gorm:"index" json:"driver_admin_id"` // Alternative field name
	RoleID      uint      `gorm:"index" json:"role_id"`
	AssignedAt  time.Time `json:"assigned_at"`
	
	// Relationships
	Role DriverAdminRole `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}
