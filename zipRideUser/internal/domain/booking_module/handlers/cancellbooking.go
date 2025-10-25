package handlers

import (
	"net/http"
	"zipride/internal/domain/booking_module/repository"

	"github.com/gin-gonic/gin"
)



//cancell booking
func CancelBookingHandler(c *gin.Context) {
	var req struct {
		BookingID uint   `json:"booking_id" binding:"required"`
		Reason    string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}	

	if err := repository.CancelBooking(req.BookingID, req.Reason); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	updatedBooking, _ := repository.GetBookingByID(req.BookingID)
	c.JSON(http.StatusOK, gin.H{"message": "booking cancelled successfully", "data": updatedBooking})
}