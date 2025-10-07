package models

import (
	"gorm.io/gorm"
)

// manager && staff && admin

type Admin struct {
	gorm.Model
	Name        string       `json:"name"`
	Email       string       `json:"email"`
	PhoneNumber string       `json:"phonenumber"`
	Password    string       `json:"password"`
	RoleID      string       `json:"role_id"`
	Role        Role         `json:"role"`
	Permissions []Permission `gorm:"many2many:admin_permissions;"`
	Block       bool         `gorm:"default:false"`
}
