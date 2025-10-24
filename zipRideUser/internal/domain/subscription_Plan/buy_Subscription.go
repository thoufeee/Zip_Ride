package subscriptionplan

import (
	"net/http"
	"strconv"
	"time"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// buy subscription
func BuySubscription(c *gin.Context) {
	user_idstr := c.Param("id")

	user_Id, _ := strconv.ParseUint(user_idstr, 10, 64)

	var input struct {
		PlanID string `json:"plan_id" binding:"required"`
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

	var existing models.UserSubscription

	if err := database.DB.Where("user_id = ? AND status = ?", user_Id, constants.Subscription_Active).
		First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "You Already Have a Active Subscription"})
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", user_Id).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "user not found"})
		return
	}

	start := time.Now()
	enddate := time.Now().AddDate(0, 0, plan.DurationDays)

	user_Sub := &models.UserSubscription{
		ID:        uuid.NewString(),
		UserID:    uint(user_Id),
		UserName:  user.FirstName + " " + user.LastName,
		PlanID:    input.PlanID,
		PlanName:  plan.PlanName,
		StartDate: start,
		EndDate:   enddate,
		Status:    constants.Subscription_Active,
	}

	if err := database.DB.Create(&user_Sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to Buy Subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "Successfuly Buyed Subscription"})
}
