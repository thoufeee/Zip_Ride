package payment

import (
	"net/http"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// verifying payment
func VerifyPayment(c *gin.Context) {
	var CallbackData map[string]interface{}

	if err := c.BindJSON(&CallbackData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid callback data"})
		return
	}

	merchantTxID := CallbackData["transactionId"].(string)
	code := CallbackData["code"].(string)

	var payment models.Payment

	if err := database.DB.First(&payment, "transaction_id = ?", merchantTxID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "payment not found"})
		return
	}

	switch code {
	case "PAYMENT_SUCCESS":
		payment.Status = constants.PaymentSuccessful
	case "PAYMENT_PENDING":
		payment.Status = constants.PaymentPending
	default:
		payment.Status = constants.PaymentFailed
	}

	if err := database.DB.Save(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to update payment status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "Payment Status Updated",
		"payment": payment,
	})
}
