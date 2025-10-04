package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"zipride/internal/driverAuth/services"
)

var authServiceLogin = services.NewAuthService()

// EmailLoginHandler verifies email/password and returns tokens or requests phone verification
func EmailLoginHandler(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := authServiceLogin.LoginWithEmail(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// OTPLoginHandler - login with OTP: client sends phone + otp; this handler verifies and issues tokens via FinalizeOtpVerification
func OTPLoginHandler(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Otp   string `json:"otp" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify OTP first
	valid, err := otpService.VerifyOTP(context.Background(), req.Phone, req.Otp)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid otp"})
		return
	}

	// finalize: mark phone verified and issue tokens if approved
	res, err := authServiceLogin.FinalizeOtpVerification(context.Background(), req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
