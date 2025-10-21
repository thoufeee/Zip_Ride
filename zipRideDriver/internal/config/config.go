package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port int
	AppEnv string

	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	TwilioAccountSID string
	TwilioAuthToken  string
	TwilioPhone      string

	SendGridAPIKey   string
	SendGridFrom     string
	SendGridFromName string

	JWTSecret          string
	AccessTokenExpiry  int
	RefreshTokenExpiry int

	DriverAdminEmail    string
	DriverAdminPassword string

	GoogleClientID string
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{}

	cfg.AppEnv = getEnv("APP_ENV", "development")
	cfg.Port = getEnvAsInt("PORT", 8080)

	cfg.DBHost = getEnv("DB_HOST", "localhost")
	cfg.DBPort = getEnvAsInt("DB_PORT", 5432)
	cfg.DBUser = strings.TrimSpace(getEnv("DB_USER", "postgres"))
	cfg.DBPassword = getEnv("DB_PASSWORD", "")
	cfg.DBName = strings.TrimSpace(getEnv("DB_NAME", "zipride_db"))

	cfg.RedisAddr = getEnv("REDIS_ADDR", "localhost:6379")
	cfg.RedisPassword = getEnv("REDIS_PASSWORD", "")
	cfg.RedisDB = getEnvAsInt("REDIS_DB", 0)

	cfg.TwilioAccountSID = strings.TrimSpace(getEnv("TWILIO_ACCOUNT_SID", ""))
	cfg.TwilioAuthToken = strings.TrimSpace(getEnv("TWILIO_AUTH_TOKEN", ""))
	cfg.TwilioPhone = strings.TrimSpace(getEnv("TWILIO_PHONE_NUMBER", ""))

	cfg.SendGridAPIKey = strings.TrimSpace(getEnv("SENDGRID_API_KEY", ""))
	cfg.SendGridFrom = strings.TrimSpace(getEnv("SENDGRID_FROM_EMAIL", ""))
	cfg.SendGridFromName = strings.TrimSpace(getEnv("SENDGRID_FROM_NAME", "ZipRide"))

	cfg.JWTSecret = getEnv("JWT_SECRET", "dev-secret-change-me")
	cfg.AccessTokenExpiry = getEnvAsInt("ACCESS_TOKEN_EXPIRY", 15)
	cfg.RefreshTokenExpiry = getEnvAsInt("REFRESH_TOKEN_EXPIRY", 10080)

	cfg.DriverAdminEmail = strings.TrimSpace(getEnv("DRIVER_ADMIN_EMAIL", "admin@zipride.local"))
	cfg.DriverAdminPassword = strings.TrimSpace(getEnv("DRIVER_ADMIN_PASSWORD", "Admin@123"))

	cfg.GoogleClientID = strings.TrimSpace(getEnv("GOOGLE_CLIENT_ID", ""))
	return cfg, nil
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(v)
	}
	return def
}

func getEnvAsInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok {
		v = strings.TrimSpace(v)
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
