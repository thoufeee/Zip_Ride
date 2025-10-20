package routes

import (
	"fmt"

	adminhandlers "zipRideDriver/internal/admin/handlers"
	adminmiddleware "zipRideDriver/internal/admin/middleware"
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

	// Static & Templates for SSR admin panel
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/**/*.html")

	r.GET("/health", func(c *gin.Context) {
		utils.Ok(c, "zipride-driver-service running", gin.H{
			"env":  cfg.AppEnv,
			"port": fmt.Sprintf("%d", cfg.Port),
		})
	})

	dh := handlers.NewDriverHandler(cfg, db, rdb, log)
	ddh := handlers.NewDriverDashboardHandler(cfg, db, rdb, log)
	veh := handlers.NewVehicleHandler(cfg, db, rdb, log)
	rh := handlers.NewRideHandler(cfg, db, rdb, log)
	helph := handlers.NewHelpHandler(cfg, db, rdb, log)
	ah := handlers.NewAdminHandler(cfg, db, rdb, log)
	wsh := handlers.NewWSHandler(cfg, db, rdb, log)
	auth := handlers.NewAuthHandler(cfg, db, rdb, log)

	// Driver auth
	driver := r.Group("/driver")
	driver.POST("/send-otp", auth.SendOTP)
	driver.POST("/verify-otp", auth.VerifyOTP)
	driver.POST("/login", auth.Login)
	driver.POST("/refresh-token", auth.RefreshToken)
	driver.POST("/logout", auth.Logout)
	driverAuth := driver.Group("")
	driverAuth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	driverAuth.POST("/onboarding", dh.Onboarding)

	// ==================== SSR ADMIN PANEL ====================
	// Session middleware for SSR admin
	panel := r.Group("/admin/panel")
	panel.Use(adminmiddleware.SSRAuthMiddleware(cfg, db, rdb, log))

	// Auth pages
	authSSR := adminhandlers.NewAuthHandler(cfg, db, rdb, log)
	panel.GET("/login", authSSR.LoginPage)
	panel.POST("/login", authSSR.Login)
	panel.GET("/logout", authSSR.Logout)

	// Protected admin panel
	panelAuth := panel.Group("")
	panelAuth.Use(adminmiddleware.RequireAdmin())

	// Dashboard
	dash := adminhandlers.NewDashboardHandler(db, log)
	panelAuth.GET("/dashboard", adminmiddleware.ACLMiddleware(db, "admin:dashboard:view"), dash.DashboardPage)

	// Drivers
	drv := adminhandlers.NewDriverHandler(db, log)
	panelAuth.GET("/drivers", adminmiddleware.ACLMiddleware(db, "admin:drivers:view"), drv.DriversPage)
	panelAuth.GET("/driver/:id", adminmiddleware.ACLMiddleware(db, "admin:drivers:view"), drv.DriverDetailPage)
	panelAuth.POST("/driver/:id/approve", adminmiddleware.ACLMiddleware(db, "admin:drivers:approve"), drv.ApproveDriver)
	panelAuth.POST("/driver/:id/reject", adminmiddleware.ACLMiddleware(db, "admin:drivers:reject"), drv.RejectDriver)
	panelAuth.POST("/driver/:id/suspend", adminmiddleware.ACLMiddleware(db, "admin:drivers:suspend"), drv.SuspendDriver)

	// Roles
	roles := adminhandlers.NewRoleHandler(db, log)
	panelAuth.GET("/roles", adminmiddleware.ACLMiddleware(db, "admin:roles:view"), roles.RolesPage)
	panelAuth.GET("/roles/new", adminmiddleware.ACLMiddleware(db, "admin:roles:edit"), roles.NewRolePage)
	panelAuth.POST("/roles", adminmiddleware.ACLMiddleware(db, "admin:roles:edit"), roles.CreateRole)
	panelAuth.GET("/roles/:id/edit", adminmiddleware.ACLMiddleware(db, "admin:roles:edit"), roles.EditRolePage)
	panelAuth.POST("/roles/:id", adminmiddleware.ACLMiddleware(db, "admin:roles:edit"), roles.UpdateRole)
	panelAuth.POST("/roles/:id/delete", adminmiddleware.ACLMiddleware(db, "admin:roles:edit"), roles.DeleteRole)

	// Admins
	admins := adminhandlers.NewAdminsHandler(db, log)
	panelAuth.GET("/admins", adminmiddleware.ACLMiddleware(db, "admin:admins:view"), admins.AdminsPage)
	panelAuth.POST("/admins", adminmiddleware.ACLMiddleware(db, "admin:admins:edit"), admins.CreateAdmin)
	panelAuth.POST("/admins/:id", adminmiddleware.ACLMiddleware(db, "admin:admins:edit"), admins.UpdateAdmin)
	panelAuth.POST("/admins/:id/delete", adminmiddleware.ACLMiddleware(db, "admin:admins:edit"), admins.DeleteAdmin)

	// Vehicles
	vehSSR := adminhandlers.NewVehiclesHandler(db, log)
	panelAuth.GET("/vehicles", adminmiddleware.ACLMiddleware(db, "admin:vehicles:view"), vehSSR.Index)
	panelAuth.GET("/vehicles/:id", adminmiddleware.ACLMiddleware(db, "admin:vehicles:view"), vehSSR.Show)
	panelAuth.POST("/vehicles/:id/verify", adminmiddleware.ACLMiddleware(db, "admin:vehicles:verify"), vehSSR.Verify)
	panelAuth.POST("/vehicles/:id/assign", adminmiddleware.ACLMiddleware(db, "admin:vehicles:assign"), vehSSR.Assign)
	panelAuth.POST("/vehicles/:id/deactivate", adminmiddleware.ACLMiddleware(db, "admin:vehicles:deactivate"), vehSSR.Deactivate)

	// Rides
	ridesSSR := adminhandlers.NewRidesHandler(db, log)
	panelAuth.GET("/rides", adminmiddleware.ACLMiddleware(db, "admin:rides:view"), ridesSSR.Index)
	panelAuth.GET("/rides/:id", adminmiddleware.ACLMiddleware(db, "admin:rides:view"), ridesSSR.Show)
	panelAuth.POST("/rides/:id/cancel", adminmiddleware.ACLMiddleware(db, "admin:rides:cancel"), ridesSSR.Cancel)

	// Earnings & Withdrawals
	earnSSR := adminhandlers.NewEarningsHandler(db, log)
	panelAuth.GET("/earnings", adminmiddleware.ACLMiddleware(db, "admin:earnings:view"), earnSSR.Index)
	panelAuth.GET("/earnings/export", adminmiddleware.ACLMiddleware(db, "admin:earnings:view"), earnSSR.ExportCSV)
	panelAuth.GET("/withdrawals", adminmiddleware.ACLMiddleware(db, "admin:earnings:view"), earnSSR.Withdrawals)
	panelAuth.POST("/withdrawals/:id/approve", adminmiddleware.ACLMiddleware(db, "admin:withdrawals:approve"), earnSSR.ApproveWithdrawal)
	panelAuth.POST("/withdrawals/:id/reject", adminmiddleware.ACLMiddleware(db, "admin:withdrawals:reject"), earnSSR.RejectWithdrawal)

	// Users (Riders)
	usersSSR := adminhandlers.NewUsersHandler(db, log)
	panelAuth.GET("/users", adminmiddleware.ACLMiddleware(db, "admin:users:view"), usersSSR.Index)
	panelAuth.POST("/users/:id/block", adminmiddleware.ACLMiddleware(db, "admin:users:block"), usersSSR.Block)
	panelAuth.POST("/users/:id/unblock", adminmiddleware.ACLMiddleware(db, "admin:users:unblock"), usersSSR.Unblock)

	// Help Center (Issues)
	helpSSR := adminhandlers.NewHelpHandler(db, log)
	panelAuth.GET("/help", adminmiddleware.ACLMiddleware(db, "admin:help:view"), helpSSR.Index)
	panelAuth.GET("/help/:id", adminmiddleware.ACLMiddleware(db, "admin:help:view"), helpSSR.Show)
	panelAuth.POST("/help/:id/reply", adminmiddleware.ACLMiddleware(db, "admin:help:reply"), helpSSR.Reply)
	panelAuth.POST("/help/:id/close", adminmiddleware.ACLMiddleware(db, "admin:help:close"), helpSSR.Close)
	// Issues aliases
	panelAuth.GET("/issues", adminmiddleware.ACLMiddleware(db, "admin:help:view"), helpSSR.Index)
	panelAuth.POST("/issues/:id/resolve", adminmiddleware.ACLMiddleware(db, "admin:help:close"), helpSSR.Close)

	// Settings
	settings := adminhandlers.NewSettingsHandler(db, log)
	panelAuth.GET("/settings", adminmiddleware.ACLMiddleware(db, "admin:settings:edit"), settings.SettingsPage)
	panelAuth.POST("/settings/password", adminmiddleware.ACLMiddleware(db, "admin:settings:edit"), settings.UpdatePassword)

	// Driver private API
	apiDriver := r.Group("/api/driver")
	apiDriver.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	apiDriver.GET("/:driverId/profile", ddh.Profile)
	apiDriver.PATCH("/:driverId/status", ddh.UpdateStatus)
	apiDriver.GET("/:driverId/earnings/summary", ddh.EarningsSummary)
	apiDriver.GET("/:driverId/earnings/trend", ddh.EarningsTrend)
	apiDriver.GET("/:driverId/rides/summary", ddh.RidesSummary)
	apiDriver.POST("/:driverId/withdraw", ddh.Withdraw)

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
