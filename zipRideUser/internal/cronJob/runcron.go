package cronjob

import (
	"log"
	"time"

	"github.com/robfig/cron"
)

// run cron for check the subscription

func RunCron() {
	c := cron.New()

	err := c.AddFunc("0 0 * * *", func() {
		log.Println("Running Subscription expiry Check", time.Now())
		CheckExpiredSubscriptions()
	})

	if err != nil {
		log.Println("failed to add cron job")
		return
	}

	c.Start()
	log.Println("subscription cron started")
}
