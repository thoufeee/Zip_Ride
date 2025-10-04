package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zipride/internal/driverAuth/services"
)

func SendDriverOtp(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "phone required"})
		return
	}

	otp, err := services.SendOtp(req.Phone)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"res": "OTP sent",
		"OTP": otp,
	})
}

func VerifyOtpHandler(c *gin.Context) {
	var data struct {
		Phone string `json:"phone" binding:"required"`
		OTP   string `json:"otp" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "phone and otp required"})
		return
	}

	ok, err := services.VerifyOtp(data.Phone, data.OTP)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	services.MarkPhoneVerified(data.Phone)
	c.JSON(http.StatusOK, gin.H{"res": "Phone verified"})
}
