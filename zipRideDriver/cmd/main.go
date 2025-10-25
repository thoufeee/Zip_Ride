package main

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"zipRideDriver/internal/config"
	"zipRideDriver/internal/database"
	"zipRideDriver/internal/routes"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config load error: %v\n", err)
		os.Exit(1)
	}

	logger, err := config.NewLogger(cfg.AppEnv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger init error: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = logger.Sync() }()

	logger.Info("starting zipride-driver-service", zap.Int("port", cfg.Port), zap.String("env", cfg.AppEnv))

	gormDB, err := database.InitPostgres(cfg, logger)
	if err != nil {
		logger.Fatal("postgres init failed", zap.Error(err))
	}

	rdb, err := database.InitRedis(cfg, logger)
	if err != nil {
		logger.Fatal("redis init failed", zap.Error(err))
	}

	if err := database.SeedDefaults(gormDB, cfg, logger); err != nil {
		logger.Fatal("seeding failed", zap.Error(err))
	}

	r := routes.SetupRouter(cfg, logger, gormDB, rdb)
	addr := fmt.Sprintf(":%d", cfg.Port)
	if err := r.Run(addr); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
}
