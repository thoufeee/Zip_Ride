package models

import (
	"gorm.io/gorm"
)

// user register
type User struct {
	gorm.Model
	GoogleID    string `json:"googleid"`
	FirstName   string `json:"firstname" binding:"required"`
	LastName    string `json:"lastname" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Gender      string `json:"gender"`
	PhoneNumber string `json:"phone"`
	Place       string `json:"place" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Role        string `json:"role"`
	Block       bool   `gorm:"default:false"`
	Isverified  bool   `json:"verified" gorm:"default:false"`
}

// google authentication

type GoogleUser struct {
	GoogleID  string
	Email     string
	FirstName string
	LastName  string
	Avatar    string
}
