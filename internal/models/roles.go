package models

// roles

type Role struct {
	ID          uint         `gorm:"primarykey" json:"id"`
	Name        string       `json:"name"`
	Permissions []Permission `gorm:"many2many:user_permissions;" json:"permissions"`
}
