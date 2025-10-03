package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zipride/internal/driverAuth/services"
)

func DriverGoogleLogin(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "token required"})
		return
	}

	token, err := services.GoogleLogin(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "login successful", "token": token})
}
