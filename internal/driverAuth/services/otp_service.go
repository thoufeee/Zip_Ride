package services

import (
	"errors"
	"fmt"
	"time"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/utils"
)

// Rate limiting constants
const (
	MaxOTPAttempts = 3
	OTPCooldown    = 60 * time.Second
	MaxDailyOTPs   = 5
)

func SendOtp(phone string) (string, error) {
	// Rate limiting checks
	if err := checkOTPRateLimit(phone); err != nil {
		return "", err
	}

	otp := utils.GeneratorOtp()
	// Save OTP in Redis with TTL
	err := utils.SaveOTP(phone, otp, constants.DriverPrefix)
	if err != nil {
		return "", err
	}

	// Send OTP via Twilio
	twilioClient := NewTwilioClient()
	message := "Your OTP code is: " + otp
	err = twilioClient.SendSMS(phone, message)
	if err != nil {
		return "", errors.New("failed to send OTP via SMS")
	}

	// Update rate limiting counters
	updateOTPRateLimit(phone)

	return otp, nil
}

func checkOTPRateLimit(phone string) error {
	// Check daily limit
	dailyKey := fmt.Sprintf("otp_daily:%s", phone)
	dailyCount, err := database.RDB.Get(database.Ctx, dailyKey).Int()
	if err == nil && dailyCount >= MaxDailyOTPs {
		return errors.New("daily OTP limit exceeded")
	}

	// Check cooldown
	cooldownKey := fmt.Sprintf("otp_cooldown:%s", phone)
	exists, err := database.RDB.Exists(database.Ctx, cooldownKey).Result()
	if err == nil && exists > 0 {
		ttl, _ := database.RDB.TTL(database.Ctx, cooldownKey).Result()
		return fmt.Errorf("please wait %d seconds before requesting another OTP", int(ttl.Seconds()))
	}

	return nil
}

func updateOTPRateLimit(phone string) {
	// Set cooldown
	cooldownKey := fmt.Sprintf("otp_cooldown:%s", phone)
	database.RDB.Set(database.Ctx, cooldownKey, "1", OTPCooldown)

	// Increment daily counter
	dailyKey := fmt.Sprintf("otp_daily:%s", phone)
	database.RDB.Incr(database.Ctx, dailyKey)
	database.RDB.Expire(database.Ctx, dailyKey, 24*time.Hour)
}

func VerifyOtp(phone, code string) (bool, error) {
	// Check attempt limit
	attemptKey := fmt.Sprintf("otp_attempts:%s", phone)
	attempts, err := database.RDB.Get(database.Ctx, attemptKey).Int()
	if err == nil && attempts >= MaxOTPAttempts {
		return false, errors.New("too many OTP attempts, please request a new OTP")
	}

	result := utils.VerifyOTP(phone, code, constants.DriverPrefix)
	if result != "valid" {
		// Increment attempt counter
		database.RDB.Incr(database.Ctx, attemptKey)
		database.RDB.Expire(database.Ctx, attemptKey, 15*time.Minute)
		return false, errors.New("invalid or expired OTP")
	}

	// Clear attempt counter on success
	database.RDB.Del(database.Ctx, attemptKey)
	return true, nil
}

func MarkPhoneVerified(phone string) {
	utils.MarkPhoneVerified(phone, "driver_verified")
}
