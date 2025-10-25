package subscriptionmanagment

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// creating subscription

func CreateSubscription(c *gin.Context) {

	var input struct {
		PlanName          string  `json:"planname"`
		Description       string  `json:"description"`
		DurationDays      int     `json:"duration_days"`
		Price             float64 `json:"price"`
		ComissionDiscount float64 `json:"comission_discount"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid input"})
		return
	}

	newsubscription := &models.SubscriptionPlan{
		ID:                uuid.New().String(),
		PlanName:          input.PlanName,
		Description:       input.Description,
		DurationDays:      input.DurationDays,
		Price:             input.Price,
		ComissionDiscount: input.ComissionDiscount,
	}

	if err := database.DB.Create(&newsubscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to create plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "succcessfully created Plan"})
}
