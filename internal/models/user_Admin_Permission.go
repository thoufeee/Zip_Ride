package models

// admin permissions

type Permission struct {
	ID   uint   `gorm:"primarykey" json:"id"`
	Name string `gorm:"unique;not null" json:"name"`
}
