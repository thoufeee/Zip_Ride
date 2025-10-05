package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"zipride/internal/driverAuth/services"
)

// SendOTP sends an OTP to the specified phone
func SendOTP(c *gin.Context) {
	var req struct {
		Phone string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if _, err := services.SendOtp(req.Phone); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent"})
}

// VerifyOTP verifies the OTP for the specified phone
func VerifyOTP(c *gin.Context) {
	var req struct {
		Phone string `json:"phone"`
		Otp   string `json:"otp"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Phone == "" || req.Otp == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ok, err := services.VerifyOtp(req.Phone, req.Otp)
	if err != nil || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	services.MarkPhoneVerified(req.Phone)
	c.JSON(http.StatusOK, gin.H{"message": "Phone verified"})
}
