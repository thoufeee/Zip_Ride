package services

import (
	"net/http"
	"zipride/database"
	"zipride/internal/middleware"
	"zipride/internal/models"
	"zipride/utils"

	"github.com/gin-gonic/gin"
)

// update user profile

func UpdateUserProfile(c *gin.Context) {
	user_id := middleware.GetUserID(c)

	if user_id == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unathorized"})
		return
	}

	var user models.User

	//    finding user
	if err := database.DB.First(&user, user_id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "user not found"})
		return
	}

	var input struct {
		FirstName   *string `json:"firstname"`
		LastName    *string `json:"lastname"`
		Email       *string `josn:"email"`
		PhoneNumber *string `json:"phonenumber"`
		Place       *string `json:"place"`
		//  Password *string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid input"})
		return
	}

	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}

	if input.LastName != nil {
		user.LastName = *input.LastName
	}

	if input.Email != nil {

		if !utils.EmailCheck(*input.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"err": "email format not valid"})
			return
		}

		if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"err": "email already registered"})
			return
		}

		user.Email = *input.Email
	}

	if input.PhoneNumber != nil {
		user.PhoneNumber = *input.PhoneNumber
	}

	if input.Place != nil {
		user.Place = *input.Place
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to update user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "successfuly updated user profile"})
}
