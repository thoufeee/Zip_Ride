package subscriptionplan

import (
	"net/http"
	"strconv"
	"time"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// buy subscription
func BuySubscription(c *gin.Context) {
	user_idstr := c.Param("id")

	user_Id, _ := strconv.ParseUint(user_idstr, 10, 64)

	var input struct {
		PlanID string `json:"plan_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid request"})
		return
	}

	var plan models.SubscriptionPlan

	if err := database.DB.Where("id = ?", input.PlanID).First(&plan).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Plan currently Not available"})
		return
	}

	user_Sub := &models.UserSubscription{
		ID:        uuid.NewString(),
		UserID:    uint(user_Id),
		PlanID:    input.PlanID,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, plan.DurationDays),
		Status:    "active",
	}

	if err := database.DB.Create(&user_Sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to Buy Subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "Successfuly Buyed Subscription"})
}
