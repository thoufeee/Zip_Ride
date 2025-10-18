package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func GenerateNumericOTP(length int) (string, error) {
	if length <= 0 {
		length = 6
	}
	const digits = "0123456789"
	var sb strings.Builder
	for i := 0; i < length; i++ {
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		sb.WriteByte(digits[nBig.Int64()])
	}
	return sb.String(), nil
}

func OTPKey(prefix, identifier string) string {
	return fmt.Sprintf("otp:%s:%s", prefix, strings.TrimSpace(identifier))
}

func StoreOTP(ctx context.Context, rdb *redis.Client, key, otp string, ttl time.Duration) error {
	return rdb.Set(ctx, key, otp, ttl).Err()
}

func VerifyOTP(ctx context.Context, rdb *redis.Client, key, otp string) (bool, error) {
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if strings.TrimSpace(val) != strings.TrimSpace(otp) {
		return false, nil
	}
	_ = rdb.Del(ctx, key).Err()
	return true, nil
}
