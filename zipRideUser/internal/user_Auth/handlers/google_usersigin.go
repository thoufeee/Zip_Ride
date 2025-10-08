package handlers

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"
	"zipride/utils"

	"github.com/gin-gonic/gin"
)

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
	token, err := utils.GenerateAccess(user.ID, user.Email, user.Role, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to create access token"})
		return
	}

	// refresh token
	refresh, err := utils.GenerateRefresh(user.ID, user.Email, user.Role, nil)

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
