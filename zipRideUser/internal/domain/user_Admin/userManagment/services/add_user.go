package services

import (
	"net/http"
	"strings"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"
	"zipride/utils"

	"github.com/gin-gonic/gin"
)

// AddUser - Create a new user securely (for admin or signup)
func AddUser(c *gin.Context) {
	var input struct {
		FirstName   string `json:"firstname" binding:"required"`
		LastName    string `json:"lastname" binding:"required"`
		Email       string `json:"email" binding:"required"`
		Gender      string `json:"gender"`
		PhoneNumber string `json:"phone" binding:"required"`
		Place       string `json:"place" binding:"required"`
		Password    string `json:"password" binding:"required"`
		Role        string `json:"role"`
		Block       bool   `json:"block"`
	}

	// Bind and validate JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials", "details": err.Error()})
		return
	}

	// Validate email format
	if !utils.EmailCheck(input.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Validate phone number
	phone, ok := utils.PhoneNumberCheck(input.PhoneNumber)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number format"})
		return
	}
	//Check if user already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	if err := database.DB.Where("phone_number = ?", phone).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Phonenumber already registered"})
		return
	}
	//Hash the password
	hashed, err := utils.GenerateHash(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to secure password"})
		return
	}
	// Create new user record
	newUser := models.User{
		FirstName:   strings.TrimSpace(input.FirstName),
		LastName:    strings.TrimSpace(input.LastName),
		Email:       input.Email,
		Gender:      input.Gender,
		PhoneNumber: phone,
		Place:       strings.TrimSpace(input.Place),
		Password:    hashed,
		Role:        constants.RoleUser,
	}

	//Save to database
	if err := database.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new user", "details": err.Error()})
		return
	}

	// Success response (exclude password)
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":        newUser.ID,
			"firstname": newUser.FirstName,
			"lastname":  newUser.LastName,
			"email":     newUser.Email,
			"phone":     newUser.PhoneNumber,
			"place":     newUser.Place,
			"gender":    newUser.Gender,
		},
	})
}
