package utils

import (
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// send email

func SendEmail(toEmail, subject, plainTextContent, htmlContent string) {
	from := mail.NewEmail("ZipRide", "ride@zipRide.com")
	to := mail.NewEmail("", toEmail)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)

	if err != nil {
		log.Println("failed to send email")
		return
	}

	log.Println("Email Sent", response.StatusCode)
}
