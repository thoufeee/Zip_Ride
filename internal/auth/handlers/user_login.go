package handlers

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"
	"zipride/utils"

	"github.com/gin-gonic/gin"
)

// user signin
func SignIn(c *gin.Context) {
	var data struct {
		PhoneNumber string `json:"phone"`
		Email       string `json:"email"`
		Password    string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "fill blanks"})
		return
	}

	var user models.User

	// check phonenumber
	if err := database.DB.Where("phone_number = ?", data.PhoneNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid phonenumber or password"})
		return
	}

	// check password
	if !utils.CheckPass(user.Password, data.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid phonenumber or password"})
		return
	}

	// creating access token
	accesstoken, err := utils.GenerateAccess(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to create access token"})
		return
	}

	// creating refresh token
	refreshtoken, err := utils.GenerateRefresh(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to create refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "Successfuly LogedIn",
		"access":  accesstoken,
		"refresh": refreshtoken,
		"role":    user.Role,
	})

}
