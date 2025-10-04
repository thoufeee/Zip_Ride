package models

import "time"

// manager && staff && admin

type Admin struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phonenumber"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	Block       bool   `gorm:"default:false"`
	CreatedAt   time.Time
	DeletedAt   time.Time
}
