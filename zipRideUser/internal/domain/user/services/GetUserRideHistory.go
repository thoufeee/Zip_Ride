package services

import (
	"net/http"
	"zipride/database"
	"zipride/internal/domain/bookingmodule/handlers"
	"zipride/internal/middleware"

	"github.com/gin-gonic/gin"
)

func GetUserRideHistory(c *gin.Context) {
	//Get user ID using helper
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	//Initialize repository
	repo := bookingmodule.NewBookingRepo(database.DB)

	//Fetch rides
	rides, err := repo.GetRidesByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ride history"})
		return
	}

	//No rides found → return empty array
	if len(rides) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
			"rides":   []any{},
		})
		return
	}

	//Success
	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"total":   len(rides),
		"rides":   rides,
	})
}

