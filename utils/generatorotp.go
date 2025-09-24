package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"
	"zipride/database"
)

// create random otp
func GeneratorOtp() string {
	otp, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		log.Fatal("failed to create otp")
	}

	return fmt.Sprintf("%04d", otp.Int64())
}

// save otp
func SaveOTP(phone, otp string) error {
	key := fmt.Sprintf("otp:%s", phone)

	return database.RDB.Set(database.Ctx, key, otp, 1*time.Minute).Err()
}

// verify otp
func VerifyOTP(code string) string {
	key := fmt.Sprintf("opt:%s", code)

	stored, err := database.RDB.Get(database.Ctx, key).Result()

	if err != nil {
		return "invalid or expired token"
	}

	database.RDB.Del(database.Ctx, key)
	return stored
}

// mark phonenumber verified
func MarkPhoneVerified(phone string) {
	database.RDB.Set(database.Ctx, "verified:"+phone, "true", 15*time.Minute)
}

// get verified phone
func GetVerifiedPhone() string {
	keys, _ := database.RDB.Keys(database.Ctx, "verified:*").Result()

	if len(keys) == 0 {
		return ""
	}

	return strings.TrimPrefix(keys[0], "verified:")
}

// clear verified phone
func ClearVerifiedPhone(phone string) {
	database.RDB.Del(database.Ctx, "verified:"+phone)
}
