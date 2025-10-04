package services

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioClient struct {
	client *twilio.RestClient
	from   string
}

func NewTwilioClient() *TwilioClient {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load env")
	}
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("TWILIO_SID"),
		Password: os.Getenv("TWILIO_AUTH_TOKEN"),
	})

	from := os.Getenv("TWILIO_PHONE_NUMBER")

	return &TwilioClient{client: client, from: from}
}

func (t *TwilioClient) SendSMS(to, message string) error {
	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(t.from)
	params.SetBody(message)

	_, err := t.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %v", err)
	}
	return nil
}
