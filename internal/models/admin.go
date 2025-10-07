package models

import "time"

// manager && staff && admin

type Admin struct {
	ID          string       `gorm:"primarykey" json:"id"`
	Name        string       `json:"name"`
	Email       string       `json:"email"`
	PhoneNumber string       `json:"phonenumber"`
	Password    string       `json:"password"`
	RoleID      string       `json:"role_id"`
	Role        Role         `json:"role"`
	Permissions []Permission `gorm:"many2many:admin_permissions;" json:"permissions"`
	Block       bool         `gorm:"default:false"`
	CreatedAt   time.Time
	DeletedAt   time.Time
}
