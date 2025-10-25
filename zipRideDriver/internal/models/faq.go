package models

import "time"

type FAQ struct {
	ID        uint      `gorm:"primaryKey"`
	Question  string    `gorm:"size:255;not null"`
	Answer    string    `gorm:"type:text;not null"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
