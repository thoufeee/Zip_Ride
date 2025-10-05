package services

import (
	"errors"
	"zipride/internal/constants"
	"zipride/utils"
)

// func SendOtp(phone string) (string, error) {
// 	otp := utils.GeneratorOtp()
// 	err := utils.SaveOTP(phone, otp, "driver_otp") // prefix specific to drivers
// 	if err != nil {
// 		return "", errors.New("failed to send otp")
// 	}
// 	// call SMS service if you want
// 	return otp, nil
// }

func SendOtp(phone string) (string, error) {
	otp := utils.GeneratorOtp()
	// Save OTP in Redis (existing logic)
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

	return otp, nil
}

func VerifyOtp(phone, code string) (bool, error) {
	result := utils.VerifyOTP(phone, code, constants.DriverPrefix)
	if result != "valid" {
		return false, errors.New("invalid or expired OTP")
	}
	return true, nil
}

func MarkPhoneVerified(phone string) {
	utils.MarkPhoneVerified(phone, "driver_verified")
}
