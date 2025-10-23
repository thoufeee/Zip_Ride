package ratingdriver

import (
	"fmt"
	"net/http"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/middleware"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

func RateDriver(c *gin.Context) {
	//Get booking_id from URL param
	bookingIdParam := c.Param("booking_id")
	var bookingID uint
	_, err := fmt.Sscan(bookingIdParam, &bookingID)
	if err != nil || bookingID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	//Get current user ID from JWT/session
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	//Fetch the booking and verify it belongs to the user
	var booking models.Booking
	if err := database.DB.First(&booking, "id = ? AND user_id = ?", bookingID, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	// //Ensure ride is completed before allowing rating
	if booking.Status != constants.StatusCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot rate a ride that is not completed"})
		return
	}

	//Bind rating + optional feedback
	var req struct {
		Rating   int    `json:"rating" binding:"required,min=1,max=5"`
		Feedback string `json:"feedback"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Update booking with rating
	booking.Rating = &req.Rating
	if req.Feedback != "" {
		booking.Feedback = &req.Feedback
	}

	// Save changes to DB
	if err := database.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save rating"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Rating submitted successfully",
		"data":    booking,
	})
}
