package models

import "time"

type AdminUser struct {
	ID           uint      `gorm:"primaryKey"`
	Name         string    `gorm:"size:100;not null"`
	Email        string    `gorm:"size:120;uniqueIndex;not null"`
	PasswordHash string    `gorm:"size:255;not null"`
	Role         string    `gorm:"size:40;default:'support';index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Role struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:60;uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Permission struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:80;uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RolePermission struct {
	ID           uint `gorm:"primaryKey"`
	RoleID       uint `gorm:"index;not null"`
	PermissionID uint `gorm:"index;not null"`
}

func (RolePermission) TableName() string { return "role_permissions" }

type UserRole struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"index;not null"`
	RoleID uint `gorm:"index;not null"`
}

func (UserRole) TableName() string { return "user_roles" }
