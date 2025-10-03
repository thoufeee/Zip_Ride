package models

import (
	"time"

	"github.com/google/uuid"
)

type Driver struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name            string    `json:"name"`
	Email           *string   `gorm:"uniqueIndex;default:null" json:"email,omitempty"`
	Phone           *string   `gorm:"uniqueIndex;default:null" json:"phone,omitempty"`
	Password        *string   `json:"password,omitempty"`
	LicenseNo       string    `json:"license_no"`
	AuthProvider    string    `json:"auth_provider"`
	IsPhoneVerified bool      `gorm:"default:false" json:"is_phone_verified"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
