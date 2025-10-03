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

	// admin login using email
	if data.Email != "" {
		// check phonenumber
		if err := database.DB.Where("email = ?", data.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid phonenumber or password"})
			return
		}
		// user login using phone
	} else if data.PhoneNumber != "" {
		if err := database.DB.Where("phone_number = ?", data.PhoneNumber).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid phonenumber or password"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"err": "email or phonenumber required"})
		return
	}

	// check password
	if !utils.CheckPass(user.Password, data.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid phonenumber or password"})
		return
	}

	// check user is blocked
	if user.Block {
		c.JSON(http.StatusBadRequest, gin.H{"err": "your account is blocked"})
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
