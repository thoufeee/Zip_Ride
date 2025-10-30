package payment

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

func PaymentRequest(c *gin.Context) {
	var req struct {
		PaymentID string `json:"payment_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid request"})
		return
	}

	var payment models.Payment

	if err := database.DB.First(&payment, "id = ?", payment.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "payment not found"})
		return
	}

	apiUrl := "https://api-preprod.phonepe.com/apis/pg-sandbox/pg/v1/pay"

	merchantID := "mock id"
	saltkey := "salt key"
	saltIndex := "1"

	phonePayLoad := map[string]interface{}{
		"merchantId":        merchantID,
		"transactionId":     payment.TransactionID,
		"amount":            payment.TotalAmount * 100,
		"redirectUrl":       "https://yourapp.com/payment/verify",
		"callbackUrl":       "https://yourapp.com/api/payment/verify",
		"paymentInstrument": map[string]string{"type": "PAY_PAGE"},
	}

	payloadBytes, err := json.Marshal(phonePayLoad)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to encode payload"})
		return
	}

	bas64payload := base64.StdEncoding.EncodeToString(payloadBytes)

	hash := sha256.Sum256([]byte(bas64payload + apiUrl + saltkey))
	xVerify := hex.EncodeToString(hash[:]) + "###" + saltIndex

	requestbody := map[string]string{
		"request": bas64payload,
	}

	bodyBytes, _ := json.Marshal(requestbody)
	reqBody := bytes.NewReader(bodyBytes)

	client := &http.Client{Timeout: 15 * time.Second}
	httpreq, _ := http.NewRequest("POST", apiUrl, reqBody)
	httpreq.Header.Set("Content-Type", "application/json")
	httpreq.Header.Set("X-VERIFY", xVerify)
	httpreq.Header.Set("X-MERCHANT-ID", merchantID)

	resp, err := client.Do(httpreq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to reach phonepe"})
		return
	}

	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(reqBody).Decode(&result)

	payment.Status = constants.PaymentPending

	if err := database.DB.Save(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to proccess payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"redirect_url": "https://api-preprod.phonepe.com/redirect/to/txnid",
		"payment_id": payment.ID,
		"redirect":   result,
	})
}
