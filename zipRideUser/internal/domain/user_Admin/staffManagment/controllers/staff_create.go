package controllers

import (
	"encoding/json"
	"net/http"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"
	"zipride/utils"

	"github.com/gin-gonic/gin"
)

// creating admin

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

	if !utils.PasswordStrength(data.Password) {
		c.JSON(http.StatusConflict, gin.H{"err": "increase password strength"})
		return
	}

	var admin models.Admin
	if err := database.DB.Where("email = ?", data.Email).First(&admin).Error; err == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "email already registered"})
		return
	}

	if err := database.DB.Where("phone_number = ?", data.PhoneNumber).First(&admin).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "phone number already registered"})
		return
	}

	hashpass, err := utils.GenerateHash(data.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to hash password"})
		return
	}

	permjson, err := json.Marshal(data.ExtraPerms)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to proccess permissions"})
		return
	}

	newStaff := models.Admin{
		Name:        data.Name,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
		Password:    hashpass,
		Role:        constants.RoleAdmin,
		Permissions: permjson,
	}

	if err := database.DB.Create(&newStaff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "creation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "Successfuly created", "role": data.Role})
}
