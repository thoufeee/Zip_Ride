package bookingmodule

import (
	"fmt"
	"net/http"
	"zipride/internal/models"
"zipride/internal/domain/bookingmodule/service"
	"github.com/gin-gonic/gin"
)

// to return fare responce for specfic vehicles
func EstimateFareHandler(c *gin.Context) {
	var req models.EstimateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	esti, err := bookingmodule.GetFareEstimates(req.PickupLat, req.PickupLong, req.DropLat, req.DropLong)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate fare"})
		return
	}
	c.JSON(http.StatusOK, esti)
}

// create booking now
func CreateBookingNowHandler(c *gin.Context) {
	userId := c.GetUint("user_id")

	var req models.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	booking, err := bookingmodule.CreateBooking(userId, req.PickupLat, req.PickupLong, req.DropLat, req.DropLong, req.VehicleType, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":          booking.ID,
		"user_id":     booking.UserID,
		"pickup_lat":  booking.PickupLat,
		"pickup_long": booking.PickupLong,
		"drop_lat":    booking.DropLat,
		"drop_long":   booking.DropLong,
		"vehicle":     booking.Vehicle,
		"fare":        booking.Fare,
		"status":      booking.Status,
		"otp":         booking.OTP,
		"created_at":  booking.CreatedAt,
		"schedule_at": booking.ScheduleAt,
	})
}

// create cooking later
func CreateLaterBookingHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req models.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	fmt.Println("Parsing schedule:", req.ScheduleDate, req.ScheduleTime)

	// Parse user-friendly date + 12-hour time
	scheduleAt, err := bookingmodule.ParseUserSchedule(req.ScheduleDate, req.ScheduleTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date or time"})
		return
	}

	booking, err := bookingmodule.CreateBooking(userID, req.PickupLat, req.PickupLong, req.DropLat, req.DropLong, req.VehicleType, &scheduleAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          booking.ID,
		"user_id":     booking.UserID,
		"pickup_lat":  booking.PickupLat,
		"pickup_long": booking.PickupLong,
		"drop_lat":    booking.DropLat,
		"drop_long":   booking.DropLong,
		"vehicle":     booking.Vehicle,
		"fare":        booking.Fare,
		"status":      booking.Status,
		"otp":         booking.OTP,
		"created_at":  booking.CreatedAt,
		"schedule_at": booking.ScheduleAt,
		 "schedule_date": booking.ScheduleDate,
        "schedule_time": booking.ScheduleTime,
	})
}
