package models

import "time"

type Rider struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:120;uniqueIndex"`
	Phone     string    `gorm:"size:20;uniqueIndex"`
	IsBlocked bool      `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
