package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	redisPkg "github.com/redis/go-redis/v9"
	"zipride/database"
	"zipride/utils"
)

// OtpService handles OTP generation, storing in Redis and sending via Twilio.
type OtpService struct {
	RedisClient *redisPkg.Client
	TTL         time.Duration
}

func NewOtpService() *OtpService {
	return &OtpService{
		RedisClient: database.RDB,
		TTL:         3 * time.Minute, // expires in 3 minutes
	}
}

func (s *OtpService) otpRedisKey(phone string) string {
	return fmt.Sprintf("otp:driver:%s", phone)
}

// SendOTP generates OTP, stores it in Redis and sends via Twilio.
func (s *OtpService) SendOTP(ctx context.Context, phone string) error {
	if phone == "" {
		return errors.New("phone required")
	}

	otp := utils.GeneratorOtp() // <-- Add parentheses to call the function

	// save OTP in redis with TTL
	key := s.otpRedisKey(phone)
	if err := s.RedisClient.Set(ctx, key, otp, s.TTL).Err(); err != nil {
		return err
	}

	// send via Twilio
	if err := s.sendSmsTwilio(phone, otp); err != nil {
		// On Twilio error, optionally delete OTP so it can't be used.
		_ = s.RedisClient.Del(ctx, key).Err()
		return err
	}

	return nil
}

// VerifyOTP checks OTP stored in Redis; if valid, deletes the key.
func (s *OtpService) VerifyOTP(ctx context.Context, phone, otp string) (bool, error) {
	if phone == "" || otp == "" {
		return false, errors.New("phone and otp required")
	}
	key := s.otpRedisKey(phone)
	stored, err := s.RedisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redisPkg.Nil {
			return false, errors.New("otp expired or not found")
		}
		return false, err
	}

	if strings.TrimSpace(stored) != strings.TrimSpace(otp) {
		return false, errors.New("invalid otp")
	}

	// valid -> delete key
	_ = s.RedisClient.Del(ctx, key).Err()
	return true, nil
}

// sendSmsTwilio - simple HTTP POST to Twilio API using account SID + token
func (s *OtpService) sendSmsTwilio(toPhone, otp string) error {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	from := os.Getenv("TWILIO_PHONE_NUMBER") // <-- fix here

	if accountSid == "" || authToken == "" || from == "" {
		return errors.New("twilio environment variables are not set")
	}

	msg := fmt.Sprintf("Your ZipRide OTP is: %s. It will expire in %d minutes.", otp, int(s.TTL.Minutes()))
	// Twilio API endpoint
	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)

	v := url.Values{}
	v.Set("To", toPhone)
	v.Set("From", from)
	v.Set("Body", msg)
	reqBody := strings.NewReader(v.Encode()) // <-- fixed here

	req, err := http.NewRequest("POST", urlStr, reqBody)
	if err != nil {
		return err
	}
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Twilio returns 201 if created
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("twilio send failed: status %d", resp.StatusCode)
	}

	return nil
}
