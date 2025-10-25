package handlers

import (
	"fmt"
	"math"
	"net/http"
	"time"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/domain/booking_module/repository"
	"zipride/internal/domain/mapservice"
	"zipride/internal/kafka"
	"zipride/internal/middleware"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// ptrTime helper to return *time.Time
func ptrTime(t time.Time) *time.Time {
	return &t
}

// -------------------- 1. EstimateBooking --------------------

<<<<<<< HEAD
// createBooking → handles both instant and scheduled bookings

func createBooking(c *gin.Context, req models.CreateBookingRequest, scheduleAt *time.Time) {
	userID := middleware.GetUserID(c) // get current logged-in user's ID from JWT/session

	// Check for duplicate booking
	// Prevents a user from creating multiple identical bookings
	isDup, err := repository.IsDuplicateBooking(userID, req.PickupLat, req.PickupLong, req.DropLat, req.DropLong, req.VehicleType, scheduleAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check duplicate"})
		return
	}
	if isDup {
		c.JSON(http.StatusConflict, gin.H{"error": "duplicate booking detected"})
		return
	}

	// Calculate distance & duration using map service
	distance, durationSec, err := mapservice.GetRouteDistance(req.PickupLat, req.PickupLong, req.DropLat, req.DropLong)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	durationMin := durationSec / 60 // convert seconds to minutes

	// Fetch vehicle fare configuration dynamically from DB
	var fareConfig models.Vehicle
	if err := database.DB.Where("vehicle_type = ?", req.VehicleType).First(&fareConfig).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle fare not found"})
		return
	}

	// Fare calculation formula: base fare + per km rate + per min rate
	totalFare := fareConfig.BaseFare + (fareConfig.PerKmRate * distance) + (fareConfig.PerMinRate * durationMin)

	now := time.Now()
	booking := models.Booking{
		UserID:     userID,
		PickupLat:  req.PickupLat,
		PickupLong: req.PickupLong,
		DropLat:    req.DropLat,
		DropLong:   req.DropLong,
		Vehicle:    req.VehicleType,
		Fare:       math.Round(totalFare*100) / 100, // round to 2 decimals
		Status:     models.StatusPending,            // initial booking status
		CreatedAt:  now,
		ScheduleAt: scheduleAt,
		ScheduleDate: func() string { // store date for easy querying/display
			if scheduleAt != nil {
				return scheduleAt.Format("2006-01-02")
			}
			return now.Format("2006-01-02")
		}(),
		ScheduleTime: func() string { // store time for display
			if scheduleAt != nil {
				return scheduleAt.Format("03:04 PM")
			}
			return now.Format("03:04 PM")
		}(),
		OTP: req.OTP, // one-time-password for ride verification
	}

	// Save booking in DB
	if err := repository.SaveBooking(&booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save booking"})
		return
	}

	// Send response to user
	c.JSON(http.StatusOK, gin.H{
		"message": func() string {
			if scheduleAt != nil {
				return "Booking scheduled successfully"
			}
			return "Booking created successfully"
		}(),
		"data": booking, // return booking details
	})
}

// -------------------- Handlers --------------------

// EstimateBooking → estimates distance, duration, and fares for vehicles
=======
>>>>>>> 7ec1a0a (Docker file)
func EstimateBooking(c *gin.Context) {
	var req models.EstimateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	distance, durationSec, err := mapservice.GetRouteDistance(req.PickupLat, req.PickupLong, req.DropLat, req.DropLong)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	durationMin := durationSec / 60
	minutes := int(durationSec) / 60

	// --- Specific vehicle type ---
	if req.VehicleType != "" {
		var fare models.Vehicle
		if err := database.DB.Where("vehicle_type = ?", req.VehicleType).First(&fare).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle type not found"})
			return
		}

		totalFare := fare.BaseFare + (fare.PerKmRate * distance) + (fare.PerMinRate * durationMin)
		c.JSON(http.StatusOK, gin.H{
			"distance":         fmt.Sprintf("%.2f km", distance),
			"duration":         fmt.Sprintf("%d min %d sec", minutes, int(durationSec)%60),
			"vehicle_type":     fare.VehicleType,
			"base_fare":        fare.BaseFare,
			"people_count":     fare.PeopleCount,
			"per_km_rate":      fare.PerKmRate,
			"per_min_rate":     fare.PerMinRate,
			"total_fare":       math.Round(totalFare*100) / 100,
			"currency":         "INR",
			"surge_multiplier": 1.0,
<<<<<<< HEAD
			"eta":              5, // dummy ETA (minutes for nearest driver)
=======
			"eta":              5,
>>>>>>> 7ec1a0a (Docker file)
		})
		return
	}

	// --- Estimate for all vehicles ---
	var fares []models.Vehicle
	if err := database.DB.Find(&fares).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch vehicle fares"})
		return
	}

	type VehicleEstimate struct {
		VehicleType string  `json:"vehicle_type"`
		TotalFare   float64 `json:"total_fare"`
		BaseFare    float64 `json:"base_fare"`
		PerKmRate   float64 `json:"per_km_rate"`
		PerMinRate  float64 `json:"per_min_rate"`
		Capacity    int     `json:"capacity"`
<<<<<<< HEAD
		ETA         int     `json:"eta"` // minutes
=======
		ETA         int     `json:"eta"`
>>>>>>> 7ec1a0a (Docker file)
	}

	results := make([]VehicleEstimate, 0, len(fares))
	for _, f := range fares {
		total := f.BaseFare + (f.PerKmRate * distance) + (f.PerMinRate * durationMin)
		results = append(results, VehicleEstimate{
			VehicleType: f.VehicleType,
			TotalFare:   math.Round(total*100) / 100,
			BaseFare:    f.BaseFare,
			PerKmRate:   f.PerKmRate,
			PerMinRate:  f.PerMinRate,
			Capacity:    f.PeopleCount,
<<<<<<< HEAD
			ETA:         5, // dummy value, can calculate from driver location
=======
			ETA:         5,
>>>>>>> 7ec1a0a (Docker file)
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"distance": fmt.Sprintf("%.2f km", distance),
		"duration": fmt.Sprintf("%d min %d sec", minutes, int(durationSec)%60),
		"vehicles": results,
	})
}

// -------------------- 2. CreateBookingNow --------------------

func CreateBookingNow(c *gin.Context) {
	var req models.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)

	// Check duplicate
	isDup, err := repository.IsDuplicateBooking(userID, req.PickupLat, req.PickupLong, req.DropLat, req.DropLong, req.VehicleType, ptrTime(time.Now()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check duplicate"})
		return
	}
	if isDup {
		c.JSON(http.StatusConflict, gin.H{"error": "duplicate booking detected"})
		return
	}

	// Distance & duration
	distance, durationSec, err := mapservice.GetRouteDistance(req.PickupLat, req.PickupLong, req.DropLat, req.DropLong)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	durationMin := durationSec / 60

	// Fare
	var fareConfig models.Vehicle
	if err := database.DB.Where("vehicle_type = ?", req.VehicleType).First(&fareConfig).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle fare not found"})
		return
	}

	totalFare := fareConfig.BaseFare + (fareConfig.PerKmRate * distance) + (fareConfig.PerMinRate * durationMin)

	now := time.Now()
	booking := models.Booking{
		UserID:       userID,
		PickupLat:    req.PickupLat,
		PickupLong:   req.PickupLong,
		DropLat:      req.DropLat,
		DropLong:     req.DropLong,
		Vehicle:      req.VehicleType,
		Fare:         math.Round(totalFare*100) / 100,
		Status:       constants.StatusPending,
		CreatedAt:    now,
		ScheduleAt:   ptrTime(now),
		ScheduleDate: now.Format("2006-01-02"),
		ScheduleTime: now.Format("03:04 PM"),
		OTP:          req.OTP,
	}

	if err := repository.SaveBooking(&booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save booking"})
		return
	}

	// Kafka
	event := models.BookingMessage{
		BookingID:  booking.ID,
		UserID:     booking.UserID,
		Vehicle:    booking.Vehicle,
		PickupLat:  booking.PickupLat,
		PickupLong: booking.PickupLong,
		DropLat:    booking.DropLat,
		DropLong:   booking.DropLong,
		Fare:       booking.Fare,
		Status:     booking.Status,
	}
	if err := kafka.Producer(event); err != nil {
		fmt.Println("⚠️ Failed to send Kafka message:", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking created successfully",
		"data":    booking,
	})
}

// -------------------- 3. CreateBookingLater --------------------

func CreateBookingLater(c *gin.Context) {
	var req models.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse schedule datetime
	if req.ScheduleAt == nil && req.ScheduleDate != "" && req.ScheduleTime != "" {
		layout := "2006-01-02 03:04 PM"
		t, err := time.Parse(layout, req.ScheduleDate+" "+req.ScheduleTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule_date or schedule_time"})
			return
		}
		req.ScheduleAt = &t
	}

	if req.ScheduleAt == nil || req.ScheduleAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "schedule_at must be a future time"})
		return
	}

	userID := middleware.GetUserID(c)

	// Check duplicate
	isDup, err := repository.IsDuplicateBooking(userID, req.PickupLat, req.PickupLong, req.DropLat, req.DropLong, req.VehicleType, req.ScheduleAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check duplicate"})
		return
	}
	if isDup {
		c.JSON(http.StatusConflict, gin.H{"error": "duplicate booking detected"})
		return
	}

	// Distance & duration
	distance, durationSec, err := mapservice.GetRouteDistance(req.PickupLat, req.PickupLong, req.DropLat, req.DropLong)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	durationMin := durationSec / 60

	// Fare
	var fareConfig models.Vehicle
	if err := database.DB.Where("vehicle_type = ?", req.VehicleType).First(&fareConfig).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle fare not found"})
		return
	}

	totalFare := fareConfig.BaseFare + (fareConfig.PerKmRate * distance) + (fareConfig.PerMinRate * durationMin)

	booking := models.Booking{
		UserID:       userID,
		PickupLat:    req.PickupLat,
		PickupLong:   req.PickupLong,
		DropLat:      req.DropLat,
		DropLong:     req.DropLong,
		Vehicle:      req.VehicleType,
		Fare:         math.Round(totalFare*100) / 100,
		Status:       constants.StatusPending,
		CreatedAt:    time.Now(),
		ScheduleAt:   req.ScheduleAt,
		ScheduleDate: req.ScheduleAt.Format("2006-01-02"),
		ScheduleTime: req.ScheduleAt.Format("03:04 PM"),
		OTP:          req.OTP,
	}
	//save booking
	if err := repository.SaveBooking(&booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking scheduled successfully",
		"data":    booking,
	})
}
