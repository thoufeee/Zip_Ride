package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"zipride/internal/driverAuth/services"
)

type PhoneVerificationHandler struct {
	OtpService *services.OtpService
}

func NewPhoneVerificationHandler(otpService *services.OtpService) *PhoneVerificationHandler {
	return &PhoneVerificationHandler{OtpService: otpService}
}

func (h *PhoneVerificationHandler) SendOTP(c *gin.Context) {
	var req struct {
		Phone string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.OtpService.SendOTP(context.Background(), req.Phone); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent"})
}

func (h *PhoneVerificationHandler) VerifyOTP(c *gin.Context) {
	var req struct {
		Phone string `json:"phone"`
		Otp   string `json:"otp"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	valid, err := h.OtpService.VerifyOTP(context.Background(), req.Phone, req.Otp)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Phone verified"})
}
