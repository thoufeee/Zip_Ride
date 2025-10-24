package subscriptionuser

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// all the details of subscribed user

func SubScribedUser(c *gin.Context) {
	var users []models.UserSubscription

	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "no subscribed users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": users})
}
