package subscriptionplan

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// get all plans

func GetAllPlans(c *gin.Context) {
	var plans []models.SubscriptionPlan

	if err := database.DB.Find(&plans).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "no plans found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": plans})
}
