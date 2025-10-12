package models

import "time"
///booking struct to store booking details
type Booking struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    UserID    uint      `json:"user_id"`
    PickupLat float64   `json:"pickup_lat"`
    PickupLong float64  `json:"pickup_long"`
    DropLat   float64   `json:"drop_lat"`
    DropLong  float64   `json:"drop_long"`
    Vehicle   string    `json:"vehicle"`
    Fare      float64   `json:"fare"`
    Status    string    `json:"status"` // pending, assigned, completed
    OTP        string    `json:"otp"`   
    CreatedAt time.Time `json:"created_at"`
    ScheduleAt  time.Time `json:"schedule_at"`  // for "later" booking
    ScheduleDate string  `json:"schedule_date"` // "2025-10-10"
	ScheduleTime string  `json:"schedule_time"` // "10:00 AM"
}

//create booking request for specific now or later
type CreateBookingRequest struct {
	PickupLat    float64 `json:"pickup_lat"`
	PickupLong   float64 `json:"pickup_long"`
	DropLat      float64 `json:"drop_lat"`
	DropLong     float64 `json:"drop_long"`
	VehicleType  string  `json:"vehicle_type"`
    Fare      float64   `json:"fare"`
    OTP        string    `json:"otp"`   
	ScheduleDate string  `json:"schedule_date"` // "2025-10-10"
	ScheduleTime string  `json:"schedule_time"` // "10:00 AM"
}

// request for the user service
type EstimateRequest struct {
	PickupLat  float64 `json:"pickup_lat"`
	PickupLong float64 `json:"pickup_long"`
	DropLat    float64 `json:"drop_lat"`
	DropLong   float64 `json:"drop_long"`
}