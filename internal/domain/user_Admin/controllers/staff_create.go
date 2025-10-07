package controllers

import (
	"net/http"

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
		c.JSON(http.StatusBadRequest, gin.H{"err": "fill blanks"})
		return
	}
}
