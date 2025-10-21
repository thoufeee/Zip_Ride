package services

import (
	"fmt"
	"net/http"
	"zipride/database"
	bookingmodule "zipride/internal/domain/bookingmodule/handlers"
	"zipride/internal/middleware"

	"github.com/gin-gonic/gin"
)

func CancelBookingHandler(c *gin.Context) {
	//booking id from param
	bookingIdparam := c.Param("id")
	var bookinId uint
	_, err := fmt.Sscanf(bookingIdparam, "%d", &bookinId)

	if err != nil || bookinId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking id"})
		return
	}
	//get user id from middleware
	userid := middleware.GetUserID(c)
	if userid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unautharised"})
		return
	}

	//cancel booking using repository
	repo := bookingmodule.NewBookingRepo(database.DB)
	err = repo.CancelBooking(bookinId, userid)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message":    "Booking cancelled successfully",
		"booking_id": bookinId,
	})

}
