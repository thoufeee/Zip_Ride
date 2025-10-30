package payment

import (
	"net/http"
	"time"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// create new payment

func CreatePaymentRequest(c *gin.Context) {

	var req struct {
		BookingID      string `json:"booking_id"`
		PaymentMethod  string `json:"payment_method"`
		PaymentGateway string `json:"payment_currency"`
		Currency       string `json:"currency"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid input"})
		return
	}

	var booking models.Booking

	if err := database.DB.First(&booking, "id = ?", booking.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "Booking not found"})
		return
	}

	var payment models.Payment

	if err := database.DB.First(&payment, "booking_id =?", booking.ID).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"err": "payment already exist for this booking"})
		return
	}

	newpayment := &models.Payment{
		ID:             uuid.New(),
		UserID:         booking.UserID,
		BookingID:      booking.ID,
		TotalAmount:    booking.Fare,
		Currency:       req.Currency,
		PaymentMethod:  req.PaymentMethod,
		PaymentGateway: req.PaymentGateway,
		Status:         constants.PaymentPending,
		TransactionID:  "TXN-" + uuid.NewString(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := database.DB.Create(&newpayment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to create payment"})
		return
	}

	newpayment.Status = constants.PaymentInitialized

	if err := database.DB.Save(&newpayment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to update payment status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "Payment Initialzied Successfuly",
		"payment": payment,
	})

}
