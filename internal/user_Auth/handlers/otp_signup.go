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

// send otp to user

func SendOtpHandler(c *gin.Context) {
	var data struct {
		Phone string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&data); err != nil || data.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Phone required"})
		return
	}

	otp := utils.GeneratorOtp()
	services.SendOtp(data.Phone, "Your OTP Is "+otp)
	utils.SaveOTP(data.Phone, otp, constants.UserPrefix)

	c.JSON(http.StatusOK, gin.H{"res": "OTP Sent"})
}

// Verify OTP

func VerifyOtpHandler(c *gin.Context) {
	var data struct {
		OTP   string `json:"code"`
		Phone string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&data); err != nil || data.OTP == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "enter OTP"})
		return
	}

	phone := utils.VerifyOTP(data.Phone, data.OTP, constants.UserPrefix)

	if phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid or expired otp"})
		return
	}

	//  storing marked verified phone numbers
	utils.MarkPhoneVerified(phone, constants.UserPrefix)

	c.JSON(http.StatusOK, gin.H{"res": "PhoneNumber verified"})

}

// register user

func RegisterUser(c *gin.Context) {
	var data struct {
		FirstName string `json:"firstname" binding:"required"`
		LastName  string `json:"lastname" binding:"required"`
		Password  string `json:"password" binding:"required"`
		Email     string `json:"email" binding:"required"`
		Gender    string `json:"gender"`
		Place     string `json:"place" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid input"})
		return
	}

	var user models.User

	if !user.Isverified {
		c.JSON(http.StatusForbidden, gin.H{"err": "phone number not verified"})
		return
	}

	// email check
	if !utils.EmailCheck(data.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"err": "email format not valid"})
		return
	}

	if err := database.DB.Where("email = ?", data.Email).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "email already taken"})
		return
	}

	if err := database.DB.Create(&models.User{}).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "user not created"})
		return
	}

	// pass check
	hash, err := utils.GenerateHash(data.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "password not hashed"})
		return
	}

	// getting verified phonenumber from redis
	phone := utils.GetVerifiedPhone(constants.UserPrefix)

	if phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to get verified phonenumber"})
		return
	}

	new := &models.User{
		FirstName:   data.FirstName,
		LastName:    data.LastName,
		PhoneNumber: phone,
		Password:    hash,
		Email:       data.Email,
		Gender:      data.Gender,
		Place:       data.Place,
	}

	if err := database.DB.Create(&new).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "account not created"})
		return
	}

	utils.ClearVerifiedPhone(phone, constants.UserPrefix)

	c.JSON(http.StatusOK, gin.H{"res": "Signup Successfuly Completed"})
}
