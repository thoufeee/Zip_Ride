package handlers

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"
	"zipride/internal/user_Auth/services"
	"zipride/utils"

	"github.com/gin-gonic/gin"
)

// forget password
func ForgetPassword(c *gin.Context) {
	var input struct {
		Phone string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credential"})
		return
	}

	// Normalize phone
	phone, ok := utils.PhoneNumberCheck(input.Phone)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number format"})
		return
	}

	// Verify phone exists
	var user models.User
	if err := database.DB.Where("phone_number = ?", phone).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number is not registered"})
		return
	}

	// Generate OTP
	otp := utils.GeneratorOtp()

	// Send OTP via Twilio in E.164 format
	twilioPhone := "+91" + phone
	services.SendOtp(twilioPhone, "Your OTP is "+otp)

	// Save OTP in Redis with prefix "forgot"
	if err := utils.SaveOTP(phone, otp, "forgot"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "OTP Sent"})
}

// verify forgot password with phone number and otp

func VerifyForgotOTP(c *gin.Context) {
	//struct to get to veryfy with user mobile number for secure
	var input struct {
		Phone string `json:"phone" binding:"required"`
		Otp   string `json:"otp" binding:"required"`
	}
	//getting data from json body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid input"})
		return
	}

	// Normalize phone
	phone, ok := utils.PhoneNumberCheck(input.Phone)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid phone number format"})
		return
	}

	// Verify OTP
	result := utils.VerifyOTP(phone, input.Otp, "forgot")
	if result != "valid" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": result})
		return
	}

	// Temporarily mark phone verified
	utils.MarkPhoneVerified(phone, "forgot")
	//sucess responce
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "OTP verified successfully"})
}

// Reset PAssword
func ResetPassword(c *gin.Context) {
	var input struct {
		Phone       string `json:"phone" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	// Normalize phone
	phone, ok := utils.PhoneNumberCheck(input.Phone)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid phone number format"})
		return
	}

	// Check if phone is verified in Redis
	verified := utils.GetVerifiedPhone("forgot")
	if verified != phone {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone number not verified"})
		return
	}

	// Check if phone exists in DB
	var user models.User
	if err := database.DB.Where("phone_number = ?", phone).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone number not found"})
		return
	}

	// Hash new password
	hash, err := utils.GenerateHash(input.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Update password in DB
	if err := database.DB.Model(&models.User{}).
		Where("phone_number = ?", phone).
		Update("password", hash).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	// Clear verified phone in Redis after successful reset
	utils.ClearVerifiedPhone(phone, "forgot")

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Password reset successfully"})
}
