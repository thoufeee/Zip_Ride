package cronjob

import (
	"fmt"
	"log"
	"time"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"
	"zipride/utils"
)

// check expired subscription
func CheckExpiredSubscriptions() {
	var activeSub []models.UserSubscription

	if err := database.DB.Where("status = ?", constants.Subscription_Active).Find(&activeSub).Error; err != nil {
		log.Println("error fetching subscription")
		return
	}

	for _, sub := range activeSub {
		if time.Now().After(sub.EndDate) {
			sub.Status = constants.Subscription_Expired

			if err := database.DB.Save(&sub).Error; err != nil {
				log.Println("failed to update Subscription")
			} else {
				log.Printf("subscription expired %s\n", sub.UserName)

				subject := "Your Subscribtion Has Expired"

				body := fmt.Sprintf("Hello %s, your subscription (%s) has expired. Please renew to continue enjoying benefits.",
					sub.UserName, sub.PlanName,
				)

				utils.SendEmail(sub.UserEmail, subject, body, "<p>"+body+"<p>")
			}
		}
	}
}
