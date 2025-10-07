package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"zipride/internal/driverAuth/services"
)

// EmailLoginHandler logs a driver in using phone and password and returns an access token
func EmailLoginHandler(c *gin.Context) {
	var req struct {
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := services.LoginDriver(req.Phone, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}
