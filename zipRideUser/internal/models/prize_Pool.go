package models

import (
	"time"

	"github.com/google/uuid"
)

// price pool

type PrizePool struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	VehicleType string    `json:"vehicle_type"`
	Commission  float64   `json:"commission"`
	BonusAmount float64   `gorm:"default:0" json:"bonusamount"`
	Active      bool      `gorm:"default:true" json:"active"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	DeletedAt   time.Time `gorm:"autoDeleteTime"`
}
