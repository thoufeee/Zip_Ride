package handlers

import (
	"net/http"
	"zipride/internal/domain/booking_module/repository"
	"zipride/internal/middleware"

	"github.com/gin-gonic/gin"
)
//users history
func GetBookingHistoryHandler(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	history, err := repository.GetUserBookingHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch booking history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "booking history fetched successfully",
		"data":    history,
	})
}
