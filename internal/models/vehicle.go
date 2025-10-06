package models

import (
	"gorm.io/gorm"
	"time"
)

type Vehicle struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	DriverID     uint           `gorm:"not null;index" json:"driver_id"`
	Driver       *Driver        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"driver,omitempty"`
	Make         string         `json:"make" binding:"required"`
	Model        string         `json:"model" binding:"required"`
	Year         int            `json:"year" binding:"required"`
	Color        string         `json:"color,omitempty"`
	Registration string         `json:"registration_number" binding:"required;unique"`
	Type         string         `json:"type,omitempty"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
