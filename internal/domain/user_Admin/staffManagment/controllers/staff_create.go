package controllers

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"
	"zipride/utils"

	"github.com/gin-gonic/gin"
)

// creating staff || manager

func CreateStaff(c *gin.Context) {

	var data struct {
		Name        string   `json:"name" binding:"required"`
		Email       string   `json:"email" binding:"required"`
		PhoneNumber string   `json:"phonenumber" binding:"required"`
		Password    string   `json:"password" binding:"required"`
		Role        string   `json:"role"`
		ExtraPerms  []string `json:"extra_permissions"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid request"})
		return
	}

	// finding role
	var role models.Role

	if err := database.DB.Preload("Permissions").Where("name = ?", data.Role).First(&role).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "role not found"})
		return
	}

	// for extra permissions
	var extraPermission []models.Permission

	if len(data.ExtraPerms) > 0 {
		if err := database.DB.Where("name IN ?", data.ExtraPerms).Find(&extraPermission).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to fetch extra permissions"})
			return
		}
	}

	allPermission := append(role.Permissions, extraPermission...)

	hashpass, err := utils.GenerateHash(data.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to hash password"})
		return
	}

	newStaff := models.Admin{
		Name:        data.Name,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
		Password:    hashpass,
		RoleID:      role.ID,
		Permissions: allPermission,
	}

	if err := database.DB.Create(&newStaff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "creation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "Successfuly created", "role": data.Role})
}
