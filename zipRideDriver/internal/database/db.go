package database

import (
	"context"
	"fmt"
	"time"

	"zipRideDriver/internal/config"
	"zipRideDriver/internal/models"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitPostgres(cfg *config.Config, log *zap.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)
	gormLogger := logger.New(
		zap.NewStdLog(log),
		logger.Config{SlowThreshold: 200 * time.Millisecond, LogLevel: logger.Info, IgnoreRecordNotFoundError: true, Colorful: false},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(
		&models.Driver{},
		&models.DriverDocument{},
		&models.Vehicle{},
		&models.Ride{},
		&models.Earning{},
		&models.Withdrawal{},
		&models.FAQ{},
		&models.HelpTicket{},
		&models.ChatSession{},
		&models.ChatMessage{},
		&models.AdminUser{},
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.UserRole{},
		&models.Rider{},
	); err != nil {
		return nil, err
	}
	return db, nil
}

func InitRedis(cfg *config.Config, log *zap.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr, Password: cfg.RedisPassword, DB: cfg.RedisDB})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Error("redis ping failed", zap.Error(err))
		return nil, err
	}
	return client, nil
}
