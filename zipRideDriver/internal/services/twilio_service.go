package services

import (
	"context"
	"fmt"
	"strings"

	"zipRideDriver/internal/config"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"go.uber.org/zap"
)

type TwilioClient struct {
	AccountSID string
	AuthToken  string
	From       string
	log        *zap.Logger
}

func NewTwilio(cfg *config.Config, log *zap.Logger) *TwilioClient {
	return &TwilioClient{
		AccountSID: strings.TrimSpace(cfg.TwilioAccountSID),
		AuthToken:  strings.TrimSpace(cfg.TwilioAuthToken),
		From:       strings.TrimSpace(cfg.TwilioPhone),
		log:        log,
	}
}

func (t *TwilioClient) SendOTP(ctx context.Context, toPhone, otp string) error {
	if t.AccountSID == "" || t.AuthToken == "" || t.From == "" {
		if t.log != nil {
			t.log.Info("twilio creds missing; skipping SMS", zap.String("to", toPhone), zap.String("otp", otp))
		}
		return nil
	}
	client := twilio.NewRestClientWithParams(twilio.ClientParams{Username: t.AccountSID, Password: t.AuthToken})
	msg := fmt.Sprintf("Your ZipRide verification code is %s", strings.TrimSpace(otp))
	params := &openapi.CreateMessageParams{}
	params.SetTo(toPhone)
	params.SetFrom(t.From)
	params.SetBody(msg)
	_, err := client.Api.CreateMessage(params)
	return err
}
