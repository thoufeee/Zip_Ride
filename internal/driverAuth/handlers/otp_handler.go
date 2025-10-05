package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"zipride/internal/driverAuth/services"
)

// SendDriverOtpHandler sends an OTP to the given phone number
func SendDriverOtpHandler(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone required"})
		return
	}

	if _, err := services.SendOtp(req.Phone); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "otp sent"})
}

// VerifyDriverOtpHandler verifies the OTP for the given phone number and issues a session token
func VerifyDriverOtpHandler(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Otp   string `json:"otp" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone and otp required"})
		return
	}

	ok, err := services.VerifyOtp(req.Phone, req.Otp)
	if err != nil || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired otp"})
		return
	}

	token, status, err := services.EnsureDriverByPhoneAndIssueToken(req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "status": status})
}
