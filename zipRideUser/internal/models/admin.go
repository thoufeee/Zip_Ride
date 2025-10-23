package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// admin

type Admin struct {
	gorm.Model
<<<<<<< HEAD
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	PhoneNumber string         `json:"phonenumber"`
	Password    string         `json:"password"`
	Role        string         `json:"role"`
	Permissions datatypes.JSON `json:"permissions" gorm:"type:json"`
	Block       bool           `gorm:"default:false"`
=======
	Name        string       `json:"name"`
	Email       string       `json:"email"`
	PhoneNumber string       `json:"phonenumber"`
	Password    string       `json:"password"`
	Role        string       `json:"role"`
	Permissions []Permission `gorm:"many2many:admin_permissions;"`
	Block       bool         `gorm:"default:false"`
>>>>>>> 2c00f30 (folders changed)
}
