package handlers

import (
	"net/http"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"
	"zipride/internal/user_Auth/services"
	"zipride/utils"

	"github.com/gin-gonic/gin"
)

// otp signin handler

func OtpSignin(c *gin.Context) {
	var data struct {
		PhoneNumber string `json:"phonenumber" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid input"})
		return
	}

	if !utils.PhoneNumberCheck(data.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"err": "inavlid phonenumber"})
	}

	var user models.User

	if err := database.DB.Where("phone = ?", data.PhoneNumber).Find(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "Phonenumber not registered"})
		return
	}

	if user.Block {
		c.JSON(http.StatusBadRequest, gin.H{"err": "your account is blocked"})
		return
	}

	// creating otp
	otp := utils.GeneratorOtp()

	// saving otp in redis
	if err := utils.SaveOTP(data.PhoneNumber, otp, constants.UserPrefix); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to store otp"})
		return
	}

	services.SendOtp(data.PhoneNumber, "Your OTp is "+otp)

	c.JSON(http.StatusOK, gin.H{"res": "OTP Sent",
		"phone": data.PhoneNumber,
	})
}

// otp verifify handler

func VerifyOTP(c *gin.Context) {
	var data struct {
		PhoneNumber string `json:"phone" binding:"required"`
		OTP         string `json:"otp" binding:"required"`
	}
}
