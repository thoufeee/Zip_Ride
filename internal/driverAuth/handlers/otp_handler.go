package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"zipride/internal/driverAuth/services"
)

var otpService = services.NewOtpService()
var authService = services.NewAuthService()

// SendDriverOtpHandler - send an OTP to phone (used for signup & login flows)
func SendDriverOtpHandler(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required,e164"` // e164 tag is illustrative; gin uses validator v9 if configured
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// fallback if validator tag e164 isn't configured
		if req.Phone == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "phone required"})
			return
		}
	}

	if err := otpService.SendOTP(context.Background(), req.Phone); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not send otp: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "otp sent"})
}

// VerifyDriverOtpHandler - verify OTP. If new signup allowed, create driver or finalize phone verification.
func VerifyDriverOtpHandler(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Otp   string `json:"otp" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone and otp required"})
		return
	}

	valid, err := otpService.VerifyOTP(context.Background(), req.Phone, req.Otp)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid otp"})
		return
	}

	// If OTP valid -> finalize verification (mark phone_verified true and issue tokens if approved)
	res, err := authService.FinalizeOtpVerification(context.Background(), req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
