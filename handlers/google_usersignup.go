package handlers

import (
	"net/http"
	"zipride/constants"
	"zipride/database"
	"zipride/models"
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

	var User models.GoogleUser

	result := database.DB.Where("google_id = ?", User.GoogleID).First(&googleuser)

	if result.Error == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "account already exists"})
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

// google user sigin

func GoogleSigin(c *gin.Context) {
	var data struct {
		Token string `json:"token"`
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

	if err := database.DB.Where("google_id = ?", googleuser.GoogleID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "account not found, please sign up"})
		return
	}

	// access token
	token, err := utils.GenerateAccess(user.ID, user.Email, user.Role)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to create access token"})
		return
	}

	// refresh token
	refresh, err := utils.GenerateRefresh(user.ID, user.Email, user.Role)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to create refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"res":     "Successfuly Loged",
		"access":  token,
		"refresh": refresh,
		"role":    user.Role,
	})
}
