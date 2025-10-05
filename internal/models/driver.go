package models

import (
	"gorm.io/gorm"
	"time"
)

type Driver struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	FirstName     string         `json:"first_name" binding:"required"`
	LastName      string         `json:"last_name" binding:"required"`
	Email         string         `gorm:"uniqueIndex" json:"email"`
	Phone         string         `gorm:"uniqueIndex" json:"phone"`
	Password      string         `json:"-"`
	IsVerified    bool           `json:"is_verified"`
	Role          string         `json:"role"`
	Status        string         `gorm:"default:pending_docs" json:"status"` // pending_docs, in_review, approved, rejected, suspended
	PhoneVerified bool           `gorm:"default:false" json:"phone_verified"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	GoogleID      *string        `gorm:"uniqueIndex;default:null" json:"google_id,omitempty"`
	Avatar        string         `json:"avatar,omitempty"`
}
