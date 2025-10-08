package models

// admin permissions

type Permission struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"unique;not null" json:"name"`
}
