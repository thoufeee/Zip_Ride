package models

import "time"

// BookingStatus constants
const (
	StatusPending   = "pending"
	StatusAssigned  = "assigned"
	StatusCompleted = "completed"
	StatusCancelled = "cancelled"
)

// Booking represents a ride booking
type Booking struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	UserID       uint       `gorm:"index" json:"user_id"`
	DriverID     *uint      `gorm:"index" json:"driver_id,omitempty"` // nullable until driver accepts
	PickupLat    float64    `json:"pickup_lat"`
	PickupLong   float64    `json:"pickup_long"`
	DropLat      float64    `json:"drop_lat"`
	DropLong     float64    `json:"drop_long"`
	Vehicle      string     `json:"vehicle"`
	Fare         float64    `json:"fare"`
	Status       string     `gorm:"index" json:"status"`
	CancelReason string     `json:"cancel_reason,omitempty"`
	OTP          string     `json:"otp,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	ScheduleAt  *time.Time `json:"schedule_at,omitempty"`
	ScheduleDate string     `json:"schedule_date,omitempty"`
	ScheduleTime string     `json:"schedule_time,omitempty"`

	// Driver info
	DriverName  string  `json:"driver_name,omitempty"`
	DriverPhone string  `json:"driver_phone,omitempty"`
	VehicleNo   string  `json:"vehicle_no,omitempty"`
	DriverLat   float64 `json:"driver_lat,omitempty"`
	DriverLong  float64 `json:"driver_long,omitempty"`

	// Ratings & feedback
	Rating   *int    `json:"rating,omitempty"`   // 1–5
	Feedback *string `json:"feedback,omitempty"` 
}

// CreateBookingRequest used in API requests
type CreateBookingRequest struct {
    PickupLat     float64    `json:"pickup_lat"`
    PickupLong    float64    `json:"pickup_long"`
    DropLat       float64    `json:"drop_lat"`
    DropLong      float64    `json:"drop_long"`
    VehicleType   string     `json:"vehicle_type"`
    OTP           string     `json:"otp"`
    ScheduleAt    *time.Time `json:"schedule_at"`
    ScheduleDate  string     `json:"schedule_date"`
    ScheduleTime  string     `json:"schedule_time"`
}

// EstimateRequest used for distance/fare estimation
type EstimateRequest struct {
	PickupLat   float64 `json:"pickup_lat" binding:"required"`
	PickupLong  float64 `json:"pickup_long" binding:"required"`
	DropLat     float64 `json:"drop_lat" binding:"required"`
	DropLong    float64 `json:"drop_long" binding:"required"`
	VehicleType string  `json:"vehicle_type,omitempty"`
}
