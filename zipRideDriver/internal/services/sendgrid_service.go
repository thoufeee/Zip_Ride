package services

import (
	"context"
	"fmt"
	"strings"

	"zipRideDriver/internal/config"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/zap"
)

type SendGridClient struct {
	APIKey    string
	From     string
	FromName string
	log       *zap.Logger
}

func NewSendGrid(cfg *config.Config, log *zap.Logger) *SendGridClient {
	return &SendGridClient{
		APIKey:    strings.TrimSpace(cfg.SendGridAPIKey),
		From:     strings.TrimSpace(cfg.SendGridFrom),
		FromName: strings.TrimSpace(cfg.SendGridFromName),
		log:      log,
	}
}

func (s *SendGridClient) SendOTP(ctx context.Context, toEmail, otp string) error {
	if s.APIKey == "" || s.From == "" {
		if s.log != nil {
			s.log.Info("sendgrid creds missing; skipping email", zap.String("to", toEmail), zap.String("otp", otp))
		}
		return nil
	}
	from := mail.NewEmail(s.FromName, s.From)
	to := mail.NewEmail("", toEmail)
	subject := "Your ZipRide verification code"
	plain := fmt.Sprintf("Your ZipRide verification code is %s", strings.TrimSpace(otp))
	msg := mail.NewSingleEmail(from, subject, to, plain, "")
	client := sendgrid.NewSendClient(s.APIKey)
	_, err := client.SendWithContext(ctx, msg)
	return err
}
