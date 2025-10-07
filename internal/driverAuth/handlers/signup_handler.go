package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"zipride/internal/driverAuth/services"
)

// EmailSignupHandler registers a driver account and returns an access token
func EmailSignupHandler(c *gin.Context) {
	var req struct {
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=6"`
		Phone     string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := services.RegisterDriver(req.FirstName, req.LastName, req.Email, req.Phone, req.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"access_token": accessToken,
		"refresh_token": refreshToken,
	})
}
