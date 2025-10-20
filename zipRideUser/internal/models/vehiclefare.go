package models

import "time"

// Vehicle model for fare management
type Vehicle struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	VehicleType string    `json:"vehicle_type"`
	BaseFare    float64   `json:"base_fare"`
	PerKmRate   float64   `json:"per_km_rate"`
	PerMinRate  float64   `json:"per_min_rate"`
	PeopleCount int       `json:"people_count"` 
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
