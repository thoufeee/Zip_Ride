package subscriptionmanagment

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// update subscription

func UpdateSubScription(c *gin.Context) {
	id := c.Param("id")

	var subscription models.SubscriptionPlan

	if err := database.DB.Where("id = ?", id).First(&subscription).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "subscription not found"})
		return
	}

	var input struct {
		PlanName          *string  `json:"planname"`
		Description       *string  `json:"description"`
		DurationDays      *int     `json:"duration_days"`
		Price             *float64 `json:"price"`
		ComissionDiscount *float64 `json:"comission_discount"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid input"})
		return
	}

	if input.PlanName != nil {
		subscription.PlanName = *input.PlanName
	}

	if input.Description != nil {
		subscription.Description = *input.Description
	}

	if input.DurationDays != nil {
		subscription.DurationDays = *input.DurationDays
	}

	if input.Price != nil {
		subscription.Price = *input.Price
	}

	if input.ComissionDiscount != nil {
		subscription.ComissionDiscount = *input.ComissionDiscount
	}

	if err := database.DB.Save(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to update subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "Successfuly Updated Subscription Plan"})
}
