package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"
	"zipride/database"
)

// create random otp
func GeneratorOtp() string {
	otp, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		log.Fatal("failed to create otp")
	}

	return fmt.Sprintf("%06d", otp.Int64())
}

// save otp
func SaveOTP(phone, otp, prefix string) error {
	key := fmt.Sprintf("%s:%s", prefix, phone)
	return database.RDB.Set(database.Ctx, key, otp, 1*time.Minute).Err()
}

// verify otp
func VerifyOTP(phone, code, prefix string) string {
	key := fmt.Sprintf("%s:%s", prefix, phone)

	stored, err := database.RDB.Get(database.Ctx, key).Result()
	if err != nil {
		return "invalid or expired token"
	}

	// Compare stored OTP with provided code
	if stored != code {
		return "invalid otp"
	}

	// Delete after successful verification to prevent reuse
	database.RDB.Del(database.Ctx, key)

	return "valid"
}

// mark phonenumber verified
func MarkPhoneVerified(phone, prefix string) {
	database.RDB.Set(database.Ctx, prefix+":verified_phone", phone, 15*time.Minute)
}

// get verified phone
func GetVerifiedPhone(prefix string) string {
	phone, err := database.RDB.Get(database.Ctx, prefix+":verified_phone").Result()

	if err != nil {
		return ""
	}

	return phone
}

// clear verified phone
func ClearVerifiedPhone(phone, prefix string) {
	database.RDB.Del(database.Ctx, prefix+":verified_phone")
}
