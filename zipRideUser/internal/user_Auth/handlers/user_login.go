package handlers

import (
	"encoding/json"
	"net/http"
	"zipride/database"
	"zipride/internal/models"
	"zipride/utils"

	"github.com/gin-gonic/gin"
)

// user signin
func SignIn(c *gin.Context) {
	var data struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid input"})
		return
	}

	//    find admin
	var admin models.Admin

	if err := database.DB.Where("email = ?", data.Email).First(&admin).Error; err == nil {

		if !utils.CheckPass(admin.Password, data.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid email or password"})
			return
		}

		if admin.Block {
			c.JSON(http.StatusForbidden, gin.H{"err": "account blocked"})
			return
		}

		var perms []string
		if len(admin.Permissions) > 0 {
			if err := json.Unmarshal(admin.Permissions, &perms); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to parse permissions"})
				return
			}

		}
		// generate access token
		access, err := utils.GenerateAccess(admin.ID, admin.Email, admin.Role, perms)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to create access token"})
			return
		}

		// generate access token
		refresh, err := utils.GenerateAccess(admin.ID, admin.Email, admin.Role, perms)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to create refresh token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"res":         "Successfuly Logged",
			"role":        admin.Role,
			"access":      access,
			"refresh":     refresh,
			"permissions": perms,
		})

		return
	}

	var user models.User

	// check user || admin exists
	if err := database.DB.Where("email = ?", data.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid email or password"})
		return
	}

	// check password
	if !utils.CheckPass(user.Password, data.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid email or password"})
		return
	}

	// check user is blocked
	if user.Block {
		c.JSON(http.StatusBadRequest, gin.H{"err": "your account is blocked"})
		return
	}

	// creating access token
	accesstoken, err := utils.GenerateAccess(user.ID, user.Email, user.Role, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to create access token"})
		return
	}

	// creating refresh token
	refreshtoken, err := utils.GenerateRefresh(user.ID, user.Email, user.Role, nil)
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
