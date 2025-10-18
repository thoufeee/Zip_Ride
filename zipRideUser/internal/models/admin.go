package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// manager && staff && admin

type Admin struct {
	gorm.Model
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	PhoneNumber string         `json:"phonenumber"`
	Password    string         `json:"password"`
	Role        string         `json:"role"`
	Permissions datatypes.JSON `json:"permissions" gorm:"type:json"`
	Block       bool           `gorm:"default:false"`
}
