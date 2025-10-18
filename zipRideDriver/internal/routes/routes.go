package routes

import (
	"fmt"

	"zipRideDriver/internal/config"
	"zipRideDriver/internal/handlers"
	"zipRideDriver/internal/middleware"
	"zipRideDriver/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SetupRouter(cfg *config.Config, log *zap.Logger, db *gorm.DB, rdb *redis.Client) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware(log))

	r.GET("/health", func(c *gin.Context) {
		utils.Ok(c, "zipride-driver-service running", gin.H{
			"env":  cfg.AppEnv,
			"port": fmt.Sprintf("%d", cfg.Port),
		})
	})

	handlers.RegisterDriverRoutes(r, cfg, db, rdb, log)
	handlers.RegisterDriverDashboardRoutes(r, cfg, db, rdb, log)
	handlers.RegisterVehicleRoutes(r, cfg, db, rdb, log)
	handlers.RegisterHelpRoutes(r, cfg, db, rdb, log)
	handlers.RegisterAdminRoutes(r, cfg, db, rdb, log)
	handlers.RegisterRideRoutes(r, cfg, db, rdb, log)
	handlers.RegisterWSRoutes(r, cfg, db, rdb, log)

	return r
}
