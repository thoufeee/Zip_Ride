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

// forget password
func ForgetPassword(c *gin.Context) {
	//struct for phone number
	var input struct {
		Phone string `json:"phone" binding:"required"`
	}
	//get data from json body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credential"})
		return
	}
	//veryfy number format
	if !utils.PhoneNumberCheck(input.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid phone number format"})
		return
	}
	//verufy phone number exist
	var user models.User
	if err := database.DB.Where("phone_number= ?", input.Phone).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Phone number is not registered"})
		return
	}

	//otp generation
	otp := utils.GeneratorOtp()
	// send otp
	services.SendOtp(input.Phone, "Your OTP Is "+otp)
	//save otp
	utils.SaveOTP(input.Phone, otp, constants.UserPrefix)

	c.JSON(http.StatusOK, gin.H{"res": "OTP Sent"})
}

// verify forgot password with phone number and otp

func VerifyForgotOTP(c *gin.Context) {
	//struct to get specific data
	var input struct {
		Phone string `json:"phone" binding:"required"`
		Otp   string `json:"otp" binding:"required"`
	}
	//read json body for struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid input"})
		return
	}
	//veryfy otp
	result := utils.VerifyOTP(input.Phone, input.Otp, "forgot")
	if result != "valid" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": result})
		return
	}

	//temporarily verified for password reset
	utils.MarkPhoneVerified(input.Phone, "forgot")
	//sucess responce
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "OTP verified successfully"})
}

//Reset PAssword

func ResetPassword(c *gin.Context) {
	//struct to create new password
	var input struct {
		Phone       string `json:"phone" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	//Read json body to get data to struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credential"})
		return
	}
	//veryfy phone number
	varyfied := utils.GetVerifiedPhone("forgot")
	if varyfied != input.Phone {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Phone number is nit veryfied"})
		return
	}
	//Password securing
	hash, err := utils.GenerateHash(input.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to secure password"})
		return
	}
	//save password on database
	if err := database.DB.Model(&models.User{}).Where("phone_number =?", input.Phone).Update("password", hash).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}
	//clear status
	utils.ClearVerifiedPhone(input.Phone, "forgot")
	c.JSON(http.StatusOK, gin.H{"message": "password updated sucessfully"})
}
