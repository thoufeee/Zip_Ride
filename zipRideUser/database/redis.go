package database

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

// connecting redis
func InitRedis() {

	addr := os.Getenv("REDIS_REST_URL")
	pass := os.Getenv("REDIS_REST_TOKEN")
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	})

	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Redis connection failed: %v", err))
	}

	fmt.Println("Redis connected")
}
