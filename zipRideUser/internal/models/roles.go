package models

// roles

type Role struct {
	ID          uint         `gorm:"primarykey"`
	Name        string       `json:"name"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}
