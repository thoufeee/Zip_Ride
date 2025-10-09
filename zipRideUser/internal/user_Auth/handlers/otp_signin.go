package handlers

import (
	"fmt"
	"net/http"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"
	"zipride/internal/user_Auth/services"
	"zipride/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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

	if err := database.DB.Where("phone_number = ?", phone).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "your account not registered"})
		return
	}

	if user.Block {
		c.JSON(http.StatusBadRequest, gin.H{"err": "your account is blocked"})
		return
	}

	// creating otp
	otp := utils.GeneratorOtp()

	// saving otp in redis
	if err := utils.SaveOTP(phone, otp, constants.UserPrefix); err != nil {
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

	phone, ok := utils.PhoneNumberCheck(data.PhoneNumber)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid phone number"})
		return
	}

	// get otp fro redis
	key := fmt.Sprintf("%s:%s", constants.UserPrefix, phone)
	storedOtp, err := database.RDB.Get(database.Ctx, key).Result()

	if err == redis.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "otp expired or invalid"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "redis error"})
		return
	}

	if data.OTP != storedOtp {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid otp"})
		return
	}

	var user models.User

	// otp is valid fetch user
	if err := database.DB.Where("phone_number = ?", phone).First(&user).Error; err != nil {
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

	// delete token after success
	database.RDB.Del(database.Ctx, key)

	c.JSON(http.StatusOK, gin.H{"res": "Successfuly Loged in",
		"access":  access,
		"refresh": refresh,
	})
}
