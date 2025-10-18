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

	// Instantiate handlers
	dh := handlers.NewDriverHandler(cfg, db, rdb, log)
	ddh := handlers.NewDriverDashboardHandler(cfg, db, rdb, log)
	veh := handlers.NewVehicleHandler(cfg, db, rdb, log)
	rh := handlers.NewRideHandler(cfg, db, rdb, log)
	helph := handlers.NewHelpHandler(cfg, db, rdb, log)
	ah := handlers.NewAdminHandler(cfg, db, rdb, log)
	wsh := handlers.NewWSHandler(cfg, db, rdb, log)

	// Driver auth
	driver := r.Group("/driver")
	driver.POST("/send-otp", dh.SendOTP)
	driver.POST("/verify-otp", dh.VerifyOTP)
	driver.POST("/login", dh.Login)
	driver.POST("/refresh-token", dh.RefreshToken)
	driver.POST("/logout", dh.Logout)
	driverAuth := driver.Group("")
	driverAuth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	driverAuth.POST("/onboarding", dh.Onboarding)

	// Driver private API
	apiDriver := r.Group("/api/driver")
	apiDriver.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	apiDriver.GET("/:driverId/profile", ddh.Profile)
	apiDriver.PATCH("/:driverId/status", ddh.UpdateStatus)
	apiDriver.GET("/:driverId/earnings/summary", ddh.EarningsSummary)
	apiDriver.GET("/:driverId/earnings/trend", ddh.EarningsTrend)
	apiDriver.GET("/:driverId/rides/summary", ddh.RidesSummary)
	apiDriver.POST("/:driverId/withdraw", ddh.Withdraw)
	apiDriver.GET("/:driverId/dashboard", ddh.Dashboard)

	// Vehicles & Documents
	apiDriver.GET("/:driverId/vehicles", veh.List)
	apiDriver.POST("/:driverId/vehicles", veh.Create)
	apiDriver.PUT("/:driverId/vehicles/:id", veh.Update)
	apiDriver.DELETE("/:driverId/vehicles/:id", veh.Delete)
	apiDriver.GET("/:driverId/documents", veh.ListDocuments)
	apiDriver.POST("/:driverId/documents", veh.CreateDocument)
	apiDriver.PATCH("/:driverId/documents/:id/status", veh.UpdateDocumentStatus)

	// Ride management
	apiDriver.POST("/:driverId/location", rh.UpdateLocation)
	apiDriver.GET("/:driverId/requests", rh.ListRequests)
	apiDriver.POST("/:driverId/rides/:id/accept", rh.AcceptRide)
	apiDriver.POST("/:driverId/rides/:id/cancel", rh.CancelRide)

	// Help Centre
	helpPub := r.Group("/api/help")
	helpPub.GET("/faqs", helph.ListFAQs)
	helpAuth := r.Group("/api/help")
	helpAuth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	helpAuth.POST("/ticket", helph.CreateTicket)
	helpAuth.POST("/report", helph.CreateReport)
	helpAuth.POST("/chat/start", helph.StartChat)

	// Admin
	admin := r.Group("/admin")
	admin.POST("/login", ah.Login)
	adminAuth := admin.Group("")
	adminAuth.Use(middleware.AuthMiddleware(cfg.JWTSecret), middleware.AdminOnly())
	adminAuth.GET("/drivers", ah.ListDrivers)
	adminAuth.POST("/driver/:id/approve", ah.ApproveDriver)
	adminAuth.POST("/driver/:id/ban", ah.BanDriver)
	adminAuth.GET("/dashboard", ah.Dashboard)

	// WebSockets
	ws := r.Group("/ws")
	ws.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	ws.GET("/driver/:driverId", wsh.DriverWS)

	return r
}
