package handlers

import (
	"errors"
	"net/http"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"
	"zipride/internal/user_Auth/services"
	"zipride/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// otp signin handler

func OtpSignin(c *gin.Context) {
	var data struct {
		PhoneNumber string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid input"})
		return
	}

	phone, ok := utils.PhoneNumberCheck(data.PhoneNumber)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"err": "inavlid phonenumber"})
	}

	var user models.User

	if err := database.DB.Where("phone_number = ?", data.PhoneNumber).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "Phonenumber not registered"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"err": "databse error"})
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
		"phone": phone,
	})
}

// otp verifify handler

func VerifyOTP(c *gin.Context) {
	var data struct {
		PhoneNumber string `json:"phone" binding:"required"`
		OTP         string `json:"otp" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "inavlid input"})
		return
	}

	// get otp fro redis
	key := "otp_" + data.PhoneNumber
	storedOtp, err := database.RDB.Get(database.Ctx, key).Result()

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "otp expired or invalid"})
		return
	}

	if data.OTP != storedOtp {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid otp"})
		return
	}

	var user models.User

	// otp is valid fetch user
	if err := database.DB.Where("phone = ?", data.PhoneNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "user not found"})
		return
	}

	// create access token
	access, err := utils.GenerateAccess(user.ID, user.Email, user.Role, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to create access token"})
		return
	}

	// create refresh token
	refresh, err := utils.GenerateRefresh(user.ID, user.Email, user.Role, nil)
	if err != nil {
		{
			c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to create refresh token"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"res": "Successfuly Loged in",
		"access":  access,
		"refresh": refresh,
	})
}
