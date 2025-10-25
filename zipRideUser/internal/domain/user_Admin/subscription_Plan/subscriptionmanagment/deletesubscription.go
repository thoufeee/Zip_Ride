package subscriptionmanagment

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// delete subscription
func DeleteSubscription(c *gin.Context) {
	id := c.Param("id")

	var subscription models.SubscriptionPlan

	if err := database.DB.Where("id = ?", id).First(&subscription).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "subscription not found"})
		return
	}

	if err := database.DB.Delete(&subscription).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to delete subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "Successfully Deleted Subscription"})
}
