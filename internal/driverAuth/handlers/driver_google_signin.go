package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"zipride/internal/driverAuth/services"
)

// GoogleSignIn handles Google id_token sign-in for drivers
func GoogleSignIn(c *gin.Context) {
	var req struct {
		IDToken string `json:"id_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, status, phoneVerified, err := services.GoogleLogin(req.IDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken, "status": status, "phone_verified": phoneVerified})
}
