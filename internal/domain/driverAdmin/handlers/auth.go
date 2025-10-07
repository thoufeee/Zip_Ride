package handlers

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"
	"zipride/utils"

	"github.com/gin-gonic/gin"
)

type adminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// DriverAdminLogin authenticates a driver admin and returns access/refresh tokens
func DriverAdminLogin(c *gin.Context) {
	var req adminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var admin models.DriverAdmin
	if err := database.DB.Where("email = ?", req.Email).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !admin.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "account disabled"})
		return
	}

	// Compare password
	if !utils.CheckPass(admin.PasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	accessToken, err := utils.GenerateAccess(admin.ID, admin.Email, "driver_admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	refreshToken, err := utils.GenerateRefresh(admin.ID, admin.Email, "driver_admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
