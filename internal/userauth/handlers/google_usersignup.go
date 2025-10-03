package handlers

import (
	"net/http"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"
	"zipride/utils"

	"github.com/gin-gonic/gin"
)

// google user signup

func GoogleSignup(c *gin.Context) {
	var data struct {
		Token       string `json:"token"`
		Gender      string `json:"gender"`
		Place       string `json:"place"`
		PhoneNumber string `json:"phonenumber"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid request"})
		return
	}

	googleuser, err := utils.VerifyGoogleToken(data.Token)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid google token"})
		return
	}

	var user models.User

	// google id check

	result := database.DB.Where("google_id = ?", user.GoogleID).First(&googleuser)

	if result.Error == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "account already exists"})
		return
	}

	// mobile number check
	if err := database.DB.Where("phone = ?", data.PhoneNumber).First(&user).Error; err == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "phone number already exists"})
		return
	}

	newuser := models.User{
		GoogleID:    googleuser.GoogleID,
		FirstName:   googleuser.FirstName,
		LastName:    googleuser.LastName,
		Email:       googleuser.Email,
		Gender:      data.Gender,
		PhoneNumber: data.PhoneNumber,
		Place:       data.Place,
		Role:        constants.RoleUser,
	}

	if err := database.DB.Create(&newuser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "account created successfuly",
		"user": newuser,
	})
}
