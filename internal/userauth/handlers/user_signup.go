package handlers

import (
	"net/http"
	"strings"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"
	"zipride/utils"

	"github.com/gin-gonic/gin"
)

// user signup

func SignUp(c *gin.Context) {
	var user struct {
		FirstName   string `json:"firstname" binding:"required"`
		LastName    string `json:"lastname" binding:"required"`
		Email       string `json:"email" binding:"required"`
		Gender      string `json:"gender"`
		PhoneNumber string `json:"phone" binding:"required"`
		Place       string `json:"place" binding:"required"`
		Password    string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "fill blanks"})
		return
	}

	var existing models.User

	// email check
	if !utils.EmailCheck(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"err": "email format not valid"})
		return
	}

	if err := database.DB.Where("email = ?", user.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"err": "email already registered"})
		return
	}

	// phonenumber check

	if !utils.PhoneNumberCheck(user.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"err": "phone number not valid"})
		return
	}

	if err := database.DB.Where("phone = ?", user.PhoneNumber).First(&existing).Error; err == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "phone number already registered"})
		return
	}

	// hash pass

	if !utils.PasswordStrength(user.Password) {
		c.JSON(http.StatusConflict, gin.H{"err": "increase password strength"})
		return
	}

	hashedpass, err := utils.GenerateHash(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to hash password"})
		return
	}

	// trim data

	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	user.Email = strings.TrimSpace(user.Email)
	user.PhoneNumber = strings.TrimSpace(user.PhoneNumber)

	new := &models.User{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Gender:      user.Gender,
		PhoneNumber: user.PhoneNumber,
		Place:       user.Place,
		Password:    hashedpass,
		Role:        constants.RoleUser,
	}

	if err := database.DB.Create(&new).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to create new user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "successfuly created"})

}
