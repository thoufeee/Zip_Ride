package bookingmodule

import (
	"fmt"
	"net/http"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// request for the user service
type EstimateRequest struct {
	PickupLat  float64 `json:"pickup_lat"`
	PickupLong float64 `json:"pickup_long"`
	DropLat    float64 `json:"drop_lat"`
	DropLong   float64 `json:"drop_long"`
}

// to return fare responce for specfic vehicles
func EstimateFareHandler(c *gin.Context) {
	var req EstimateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	esti, err := GetFareEstimates(req.PickupLat, req.PickupLong, req.DropLat, req.DropLong)
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

	booking, err := CreateBooking(userId, req.PickupLat, req.PickupLong, req.DropLat, req.DropLong, req.VehicleType, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}
	c.JSON(http.StatusOK, booking)
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
	scheduleAt, err := ParseUserSchedule(req.ScheduleDate, req.ScheduleTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date or time"})
		return
	}

	booking, err := CreateBooking(userID, req.PickupLat, req.PickupLong, req.DropLat, req.DropLong, req.VehicleType, &scheduleAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"booking": booking})
}

