package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"zipride/internal/driverAuth/services"
	"zipride/internal/models" // Update this import path if Vehicle is elsewhere
)

var authSvcProfile = services.NewAuthService()

// CompleteProfileHandler - driver completes profile (license, vehicle, etc.)
func CompleteProfileHandler(c *gin.Context) {
	var req struct {
		DriverID      uint           `json:"driver_id" binding:"required"` // Use uint if your model uses uint
		FullName      string         `json:"full_name" binding:"required"`
		LicenseNumber string         `json:"license_number" binding:"required"`
		Vehicle       models.Vehicle `json:"vehicle" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := authSvcProfile.CompleteProfile(context.Background(), req.DriverID, req.FullName, req.LicenseNumber, req.Vehicle); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile submitted; awaiting admin approval", "status": "pending"})
}
