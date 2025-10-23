// package routes

// import (
// 	"html/template"
// 	"net/http"
// 	"path/filepath"

// 	adminhandlers "zipRideDriver/internal/admin/handlers"
// 	adminmiddleware "zipRideDriver/internal/admin/middleware"
// 	"zipRideDriver/internal/config"
// 	"zipRideDriver/internal/handlers"
// 	"zipRideDriver/internal/middleware"

// 	"github.com/gin-gonic/gin"
// 	"github.com/redis/go-redis/v9"
// 	"go.uber.org/zap"
// 	"gorm.io/gorm"
// )

// func SetupRouter(cfg *config.Config, log *zap.Logger, db *gorm.DB, rdb *redis.Client) *gin.Engine {
// 	r := gin.New()
// 	r.Use(gin.Recovery())
// 	r.Use(middleware.LoggerMiddleware(log))

// 	// Static & Templates for SSR admin panel
// 	r.Static("/static", "./static")
// 	// Ensure base is parsed first so its default blocks don't overwrite page blocks
// 	tpl := template.New("")
// 	tpl = template.Must(tpl.ParseFiles("templates/layouts/base.html"))
// 	adminUIEnabled := false
// 	if files, _ := filepath.Glob("templates/admin/*.html"); len(files) > 0 {
// 		tpl = template.Must(tpl.ParseFiles(files...))
// 		adminUIEnabled = true
// 	}
// 	if files, _ := filepath.Glob("templates/admin/*/*.html"); len(files) > 0 {
// 		tpl = template.Must(tpl.ParseFiles(files...))
// 		adminUIEnabled = true
// 	}
// 	r.SetHTMLTemplate(tpl)

// 	// r.GET("/health", func(c *gin.Context) {
// 	// 	utils.Ok(c, "zipride-driver-service running", gin.H{
// 	// 		"env":  cfg.AppEnv,
// 	// 		"port": fmt.Sprintf("%d", cfg.Port),
// 	// 	})
// 	// })

// 	dh := handlers.NewDriverHandler(cfg, db, rdb, log)
// 	ddh := handlers.NewDriverDashboardHandler(cfg, db, rdb, log)
// 	veh := handlers.NewVehicleHandler(cfg, db, rdb, log)
// 	rh := handlers.NewRideHandler(cfg, db, rdb, log)
// 	helph := handlers.NewHelpHandler(cfg, db, rdb, log)
// 	ah := handlers.NewAdminHandler(cfg, db, rdb, log)
// 	wsh := handlers.NewWSHandler(cfg, db, rdb, log)
// 	auth := handlers.NewAuthHandler(cfg, db, rdb, log)
// 	h := handlers.NewPublicHandler(cfg, db, rdb, log)
// 	regHandler := handlers.NewDriverRegistrationHandler(db, log)
// 	apiHandler := handlers.NewAPIHandler(cfg, db, rdb, log)

// 	// Driver auth
// 	driver := r.Group("/driver")
// 	driver.POST("/send-otp", auth.SendOTP)
// 	// Public API
// 	r.POST("/register", h.RegisterDriver)
// 	r.POST("/login", h.LoginDriver)
// 	r.GET("/health", h.Health)

// 	// Driver Registration (Public)
// 	r.POST("/api/driver/register", regHandler.RegisterDriver)
// 	r.GET("/api/driver/registration-status/:email", regHandler.CheckRegistrationStatus)
// 	driver.POST("/verify-otp", auth.VerifyOTP)
// 	driver.POST("/login", auth.Login)
// 	driver.POST("/refresh-token", auth.RefreshToken)
// 	driver.POST("/logout", auth.Logout)
// 	driverAuth := driver.Group("")
// 	driverAuth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
// 	driverAuth.POST("/onboarding", dh.Onboarding)

// 	// ==================== SSR ADMIN PANEL ====================
// 	// Session middleware for SSR admin
// 	panel := r.Group("/admin/panel")
// 	panel.Use(adminmiddleware.SSRAuthMiddleware(cfg, db, rdb, log))
// 	if adminUIEnabled {
// 		panel.GET("/", func(c *gin.Context) { c.Redirect(http.StatusSeeOther, "/admin/panel/dashboard") })
// 		panel.GET("", func(c *gin.Context) { c.Redirect(http.StatusSeeOther, "/admin/panel/dashboard") })

// 		// Auth pages
// 		authSSR := adminhandlers.NewAuthHandler(cfg, db, rdb, log)
// 		panel.GET("/login", authSSR.LoginPage)
// 		panel.POST("/login", authSSR.Login)
// 		panel.GET("/logout", authSSR.Logout)

// 		// Protected admin panel
// 		panelAuth := panel.Group("")
// 		panelAuth.Use(adminmiddleware.RequireAdmin())

// 		// Dashboard
// 		dash := adminhandlers.NewDashboardHandler(db, log)
// 		panelAuth.GET("/dashboard", adminmiddleware.ACLMiddleware(db, "admin:dashboard:view"), dash.DashboardPage)

// 		// Drivers
// 		drv := adminhandlers.NewDriverHandler(db, log)
// 		panelAuth.GET("/drivers", adminmiddleware.ACLMiddleware(db, "admin:drivers:view"), drv.DriversPage)
// 		panelAuth.GET("/drivers/pending", adminmiddleware.ACLMiddleware(db, "admin:drivers:view"), drv.PendingApprovals)
// 		panelAuth.GET("/driver/:id", adminmiddleware.ACLMiddleware(db, "admin:drivers:view"), drv.DriverDetailPage)
// 		panelAuth.POST("/driver/:id/approve", adminmiddleware.ACLMiddleware(db, "admin:drivers:approve"), drv.ApproveDriver)
// 		panelAuth.POST("/driver/:id/reject", adminmiddleware.ACLMiddleware(db, "admin:drivers:reject"), drv.RejectDriver)
// 		panelAuth.POST("/driver/:id/suspend", adminmiddleware.ACLMiddleware(db, "admin:drivers:suspend"), drv.SuspendDriver)

// 		// Roles
// 		roles := adminhandlers.NewRoleHandler(db, log)
// 		panelAuth.GET("/roles", adminmiddleware.ACLMiddleware(db, "admin:roles:view"), roles.RolesPage)
// 		panelAuth.GET("/roles/new", adminmiddleware.ACLMiddleware(db, "admin:roles:edit"), roles.NewRolePage)
// 		panelAuth.POST("/roles", adminmiddleware.ACLMiddleware(db, "admin:roles:edit"), roles.CreateRole)
// 		panelAuth.GET("/roles/:id/edit", adminmiddleware.ACLMiddleware(db, "admin:roles:edit"), roles.EditRolePage)
// 		panelAuth.POST("/roles/:id", adminmiddleware.ACLMiddleware(db, "admin:roles:edit"), roles.UpdateRole)
// 		panelAuth.POST("/roles/:id/delete", adminmiddleware.ACLMiddleware(db, "admin:roles:edit"), roles.DeleteRole)

// 		// Admins
// 		admins := adminhandlers.NewAdminsHandler(db, log)
// 		panelAuth.GET("/admins", adminmiddleware.ACLMiddleware(db, "admin:admins:view"), admins.AdminsPage)
// 		panelAuth.POST("/admins", adminmiddleware.ACLMiddleware(db, "admin:admins:edit"), admins.CreateAdmin)
// 		panelAuth.POST("/admins/:id", adminmiddleware.ACLMiddleware(db, "admin:admins:edit"), admins.UpdateAdmin)
// 		panelAuth.POST("/admins/:id/delete", adminmiddleware.ACLMiddleware(db, "admin:admins:edit"), admins.DeleteAdmin)

// 		// Vehicles
// 		vehSSR := adminhandlers.NewVehiclesHandler(db, log)
// 		panelAuth.GET("/vehicles", adminmiddleware.ACLMiddleware(db, "admin:vehicles:view"), vehSSR.Index)
// 		panelAuth.GET("/vehicles/:id", adminmiddleware.ACLMiddleware(db, "admin:vehicles:view"), vehSSR.Show)
// 		panelAuth.POST("/vehicles/:id/verify", adminmiddleware.ACLMiddleware(db, "admin:vehicles:verify"), vehSSR.Verify)
// 		panelAuth.POST("/vehicles/:id/assign", adminmiddleware.ACLMiddleware(db, "admin:vehicles:assign"), vehSSR.Assign)
// 		panelAuth.POST("/vehicles/:id/deactivate", adminmiddleware.ACLMiddleware(db, "admin:vehicles:deactivate"), vehSSR.Deactivate)

// 		// Rides
// 		ridesSSR := adminhandlers.NewRidesHandler(db, log)
// 		panelAuth.GET("/rides", adminmiddleware.ACLMiddleware(db, "admin:rides:view"), ridesSSR.Index)
// 		panelAuth.GET("/rides/:id", adminmiddleware.ACLMiddleware(db, "admin:rides:view"), ridesSSR.Show)
// 		panelAuth.POST("/rides/:id/cancel", adminmiddleware.ACLMiddleware(db, "admin:rides:cancel"), ridesSSR.Cancel)

// 		// Earnings & Withdrawals
// 		earnSSR := adminhandlers.NewEarningsHandler(db, log)
// 		panelAuth.GET("/earnings", adminmiddleware.ACLMiddleware(db, "admin:earnings:view"), earnSSR.Index)
// 		panelAuth.GET("/earnings/export", adminmiddleware.ACLMiddleware(db, "admin:earnings:view"), earnSSR.ExportCSV)
// 		panelAuth.GET("/withdrawals", adminmiddleware.ACLMiddleware(db, "admin:earnings:view"), earnSSR.Withdrawals)
// 		panelAuth.POST("/withdrawals/:id/approve", adminmiddleware.ACLMiddleware(db, "admin:withdrawals:approve"), earnSSR.ApproveWithdrawal)
// 		panelAuth.POST("/withdrawals/:id/reject", adminmiddleware.ACLMiddleware(db, "admin:withdrawals:reject"), earnSSR.RejectWithdrawal)

// 		// Users (Riders)
// 		usersSSR := adminhandlers.NewUsersHandler(db, log)
// 		panelAuth.GET("/users", adminmiddleware.ACLMiddleware(db, "admin:users:view"), usersSSR.Index)
// 		panelAuth.POST("/users/:id/block", adminmiddleware.ACLMiddleware(db, "admin:users:block"), usersSSR.Block)
// 		panelAuth.POST("/users/:id/unblock", adminmiddleware.ACLMiddleware(db, "admin:users:unblock"), usersSSR.Unblock)

// 		// Help Center (Issues)
// 		helpSSR := adminhandlers.NewHelpHandler(db, log)
// 		panelAuth.GET("/help", adminmiddleware.ACLMiddleware(db, "admin:help:view"), helpSSR.Index)
// 		panelAuth.GET("/help/:id", adminmiddleware.ACLMiddleware(db, "admin:help:view"), helpSSR.Show)
// 		panelAuth.POST("/help/:id/reply", adminmiddleware.ACLMiddleware(db, "admin:help:reply"), helpSSR.Reply)
// 		panelAuth.POST("/help/:id/close", adminmiddleware.ACLMiddleware(db, "admin:help:close"), helpSSR.Close)
// 		// Issues aliases
// 		panelAuth.GET("/issues", adminmiddleware.ACLMiddleware(db, "admin:help:view"), helpSSR.Index)
// 		panelAuth.POST("/issues/:id/resolve", adminmiddleware.ACLMiddleware(db, "admin:help:close"), helpSSR.Close)

// 		// Settings
// 		settings := adminhandlers.NewSettingsHandler(db, log)
// 		panelAuth.GET("/settings", adminmiddleware.ACLMiddleware(db, "admin:settings:edit"), settings.SettingsPage)
// 		panelAuth.POST("/settings/password", adminmiddleware.ACLMiddleware(db, "admin:settings:edit"), settings.UpdatePassword)
// 	} else {
// 		// Fallback when admin UI is removed; prevents "template is undefined" errors
// 		panel.GET("/*any", func(c *gin.Context) {
// 			c.String(http.StatusServiceUnavailable, "Admin panel UI is not installed. Implement templates to enable it.")
// 		})
// 	}

// 	// Driver private API
// 	apiDriver := r.Group("/api/driver")
// 	apiDriver.Use(middleware.AuthMiddleware(cfg.JWTSecret))
// 	apiDriver.GET("/:driverId/profile", ddh.Profile)
// 	apiDriver.PATCH("/:driverId/status", ddh.UpdateStatus)
// 	apiDriver.GET("/:driverId/earnings/summary", ddh.EarningsSummary)
// 	apiDriver.GET("/:driverId/earnings/trend", ddh.EarningsTrend)
// 	apiDriver.GET("/:driverId/rides/summary", ddh.RidesSummary)
// 	apiDriver.POST("/:driverId/withdraw", ddh.Withdraw)

// 	// Vehicles & Documents
// 	apiDriver.GET("/:driverId/vehicles", veh.List)
// 	apiDriver.POST("/:driverId/vehicles", veh.Create)
// 	apiDriver.PUT("/:driverId/vehicles/:id", veh.Update)
// 	apiDriver.DELETE("/:driverId/vehicles/:id", veh.Delete)
// 	apiDriver.GET("/:driverId/documents", veh.ListDocuments)
// 	apiDriver.POST("/:driverId/documents", veh.CreateDocument)
// 	apiDriver.PATCH("/:driverId/documents/:id/status", veh.UpdateDocumentStatus)

// 	// Ride management
// 	apiDriver.POST("/:driverId/location", rh.UpdateLocation)
// 	apiDriver.GET("/:driverId/requests", rh.ListRequests)
// 	apiDriver.POST("/:driverId/rides/:id/accept", rh.AcceptRide)
// 	apiDriver.POST("/:driverId/rides/:id/cancel", rh.CancelRide)

// 	// Help Centre
// 	helpPub := r.Group("/api/help")
// 	helpPub.GET("/faqs", helph.ListFAQs)
// 	helpAuth := r.Group("/api/help")
// 	helpAuth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
// 	helpAuth.POST("/ticket", helph.CreateTicket)
// 	helpAuth.POST("/report", helph.CreateReport)
// 	helpAuth.POST("/chat/start", apiHandler.StartLiveChat)

// 	// General API endpoints (require JWT)
// 	api := r.Group("/api")
// 	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))

// 	// Vehicle endpoints (without driverId)
// 	api.GET("/vehicles", apiHandler.GetAllVehicles)
// 	api.POST("/vehicles", apiHandler.CreateVehicle)
// 	api.PUT("/vehicles/:id", apiHandler.UpdateVehicle)

// 	// Document endpoints (without driverId)
// 	api.GET("/documents", apiHandler.GetAllDocuments)
// 	api.PUT("/documents/:id", apiHandler.UpdateDocument)

// 	// Ride management endpoints
// 	api.POST("/ride/accept", apiHandler.AcceptRide)
// 	api.POST("/ride/cancel", apiHandler.CancelRide)

// 	// Chat endpoints
// 	api.POST("/chat/send", apiHandler.SendChatMessage)

// 	// Admin Routes (supports both cookie and JWT auth)
// 	admin := r.Group("/admin")
// 	admin.GET("/login", func(c *gin.Context) {
// 		c.HTML(http.StatusOK, "admin/login.html", nil)
// 	})
// 	admin.POST("/login", ah.Login) // Enhanced to support both cookie and JWT
// 	admin.GET("/logout", func(c *gin.Context) {
// 		c.SetCookie("admin_token", "", -1, "/", "", false, true)
// 		c.Redirect(http.StatusSeeOther, "/admin/login")
// 	})

// 	// Protected Admin Routes (supports both auth methods)
// 	adminAuth := admin.Group("")
// 	adminAuth.Use(func(c *gin.Context) {
// 		// Try cookie auth first
// 		if token, err := c.Cookie("admin_token"); err == nil && token != "" {
// 			c.Next()
// 			return
// 		}
// 		// Fall back to JWT auth
// 		middleware.AuthMiddleware(cfg.JWTSecret)(c)
// 		if c.IsAborted() {
// 			return
// 		}
// 		middleware.AdminOnly()(c)
// 	})
// 	adminAuth.GET("/dashboard", ah.Dashboard)
// 	adminAuth.GET("/drivers", ah.ListDrivers)
// 	adminAuth.GET("/driver/:id/approve", ah.ApproveDriver)  // HTML link support
// 	adminAuth.POST("/driver/:id/approve", ah.ApproveDriver) // API support
// 	adminAuth.GET("/driver/:id/suspend", ah.SuspendDriver)  // HTML link support
// 	adminAuth.POST("/driver/:id/ban", ah.BanDriver)
// 	adminAuth.GET("/rides", ah.ListRides)                 // Ride management
// 	adminAuth.GET("/ride/:id", ah.ShowRide)               // Ride details
// 	adminAuth.GET("/ride/:id/cancel", ah.CancelRide)      // Cancel ride
// 	adminAuth.GET("/earnings", ah.ListEarnings)           // Earnings management
// 	adminAuth.GET("/earnings/export", ah.ExportEarnings)  // Export earnings CSV
// 	adminAuth.GET("/reports", ah.ListReports)             // Reports & Analytics
// 	adminAuth.GET("/reports/export", ah.ExportReports)    // Export platform report CSV
// 	adminAuth.GET("/settings", ah.AdminSettings)          // Admin settings
// 	adminAuth.POST("/change-password", ah.ChangePassword) // Change password

// 	// WebSockets
// 	ws := r.Group("/ws")
// 	ws.Use(middleware.AuthMiddleware(cfg.JWTSecret))
// 	ws.GET("/driver/:driverId", wsh.DriverWS)

//		return r
//	}
package routes

import (
	"html/template"
	"net/http"
	"path/filepath"

	adminhandlers "zipRideDriver/internal/admin/handlers"
	adminmiddleware "zipRideDriver/internal/admin/middleware"
	"zipRideDriver/internal/config"
	"zipRideDriver/internal/handlers"
	"zipRideDriver/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SetupRouter
func SetupRouter(cfg *config.Config, log *zap.Logger, db *gorm.DB, rdb *redis.Client) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), middleware.LoggerMiddleware(log))

	r.Static("/static", "./static")
	tpl := template.New("")
	tpl = template.Must(tpl.ParseFiles("templates/layouts/base.html"))

	adminUIEnabled := false
	patterns := []string{
		"templates/admin/*.html",
		"templates/admin/*/*.html",
		"templates/admin/*/*/*.html",
	}

	for _, pattern := range patterns {
		files, err := filepath.Glob(pattern)
		if err != nil {
			log.Warn("failed to glob admin templates", zap.String("pattern", pattern), zap.Error(err))
			continue
		}
		if len(files) == 0 {
			continue
		}
		tpl = template.Must(tpl.ParseFiles(files...))
		adminUIEnabled = true
	}
	r.SetHTMLTemplate(tpl)

	// Hanlder Initialization

	authHandler := handlers.NewAuthHandler(cfg, db, rdb, log)
	driverHandler := handlers.NewDriverHandler(cfg, db, rdb, log)
	dashboardHandler := handlers.NewDriverDashboardHandler(cfg, db, rdb, log)
	vehicleHandler := handlers.NewVehicleHandler(cfg, db, rdb, log)
	rideHandler := handlers.NewRideHandler(cfg, db, rdb, log)
	helpHandler := handlers.NewHelpHandler(cfg, db, rdb, log)
	regHandler := handlers.NewDriverRegistrationHandler(db, log)
	apiHandler := handlers.NewAPIHandler(cfg, db, rdb, log)
	wsHandler := handlers.NewWSHandler(cfg, db, rdb, log)
	publicHandler := handlers.NewPublicHandler(cfg, db, rdb, log)

	// Public routes

	public := r.Group("/api/driver")
	{
		public.POST("/register", regHandler.RegisterDriver)
		public.GET("/registration-status/:email", regHandler.CheckRegistrationStatus)

		public.POST("/auth/send-otp", authHandler.SendOTP)
		public.POST("/auth/verify-otp", authHandler.VerifyOTP)
		public.POST("/auth/login", authHandler.Login)
		public.POST("/auth/logout", authHandler.Logout)
	}

	// Health Check
	r.GET("/health", publicHandler.Health)

	// Driver routes

	driver := r.Group("/api/driver")
	driver.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		driver.POST("/onboarding", driverHandler.Onboarding)
		driver.GET("/:driverId/profile", dashboardHandler.Profile)
		driver.PATCH("/:driverId/status", dashboardHandler.UpdateStatus)

		// Earnings
		driver.GET("/:driverId/earnings/summary", dashboardHandler.EarningsSummary)
		driver.GET("/:driverId/earnings/trend", dashboardHandler.EarningsTrend)
		driver.POST("/:driverId/withdraw", dashboardHandler.Withdraw)

		// Rides
		driver.GET("/:driverId/rides/summary", dashboardHandler.RidesSummary)
		driver.POST("/:driverId/location", rideHandler.UpdateLocation)
		driver.GET("/:driverId/requests", rideHandler.ListRequests)
		driver.POST("/:driverId/rides/:id/accept", rideHandler.AcceptRide)
		driver.POST("/:driverId/rides/:id/cancel", rideHandler.CancelRide)

		// Vehicles
		driver.GET("/:driverId/vehicles", vehicleHandler.List)
		driver.POST("/:driverId/vehicles", vehicleHandler.Create)
		driver.PUT("/:driverId/vehicles/:id", vehicleHandler.Update)
		driver.DELETE("/:driverId/vehicles/:id", vehicleHandler.Delete)

		// Documents
		driver.GET("/:driverId/documents", vehicleHandler.ListDocuments)
		driver.POST("/:driverId/documents", vehicleHandler.CreateDocument)
		driver.PATCH("/:driverId/documents/:id/status", vehicleHandler.UpdateDocumentStatus)
	}

	// Help center
	help := r.Group("/api/help")
	{
		help.GET("/faqs", helpHandler.ListFAQs)
	}

	helpAuth := help.Group("")
	helpAuth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		helpAuth.POST("/ticket", helpHandler.CreateTicket)
		helpAuth.POST("/report", helpHandler.CreateReport)
		helpAuth.POST("/chat/start", apiHandler.StartLiveChat)
	}

	// websockect

	ws := r.Group("/ws")
	ws.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		ws.GET("/driver/:driverId", wsHandler.DriverWS)
	}

	// admin panel routes

	panel := r.Group("/admin/panel")
	panel.Use(adminmiddleware.SSRAuthMiddleware(cfg, db, rdb, log))

	if adminUIEnabled {
		panel.GET("/", func(c *gin.Context) { c.Redirect(http.StatusSeeOther, "/admin/panel/dashboard") })
		panel.GET("/login", adminhandlers.NewAuthHandler(cfg, db, rdb, log).LoginPage)
		panel.POST("/login", adminhandlers.NewAuthHandler(cfg, db, rdb, log).Login)
		panel.GET("/logout", adminhandlers.NewAuthHandler(cfg, db, rdb, log).Logout)

		// Authenticated Admin Section
		panelAuth := panel.Group("")
		panelAuth.Use(adminmiddleware.RequireAdmin())

		// Dashboard
		dashboard := adminhandlers.NewDashboardHandler(db, log)
		panelAuth.GET("/dashboard", adminmiddleware.ACLMiddleware(db, "admin:dashboard:view"), dashboard.DashboardPage)

		// Drivers Management
		driverAdmin := adminhandlers.NewDriverHandler(db, log)
		panelAuth.GET("/drivers", adminmiddleware.ACLMiddleware(db, "admin:drivers:view"), driverAdmin.DriversPage)
		panelAuth.GET("/drivers/pending", adminmiddleware.ACLMiddleware(db, "admin:drivers:view"), driverAdmin.PendingApprovals)
		panelAuth.GET("/driver/:id", adminmiddleware.ACLMiddleware(db, "admin:drivers:view"), driverAdmin.DriverDetailPage)
		panelAuth.POST("/driver/:id/approve", adminmiddleware.ACLMiddleware(db, "admin:drivers:approve"), driverAdmin.ApproveDriver)
		panelAuth.POST("/driver/:id/reject", adminmiddleware.ACLMiddleware(db, "admin:drivers:reject"), driverAdmin.RejectDriver)
		panelAuth.POST("/driver/:id/suspend", adminmiddleware.ACLMiddleware(db, "admin:drivers:suspend"), driverAdmin.SuspendDriver)

		// Roles & Admins
		roleAdmin := adminhandlers.NewRoleHandler(db, log)
		admins := adminhandlers.NewAdminsHandler(db, log)

		panelAuth.GET("/roles", adminmiddleware.ACLMiddleware(db, "admin:roles:view"), roleAdmin.RolesPage)
		panelAuth.POST("/roles", adminmiddleware.ACLMiddleware(db, "admin:roles:edit"), roleAdmin.CreateRole)
		panelAuth.GET("/admins", adminmiddleware.ACLMiddleware(db, "admin:admins:view"), admins.AdminsPage)
		panelAuth.POST("/admins", adminmiddleware.ACLMiddleware(db, "admin:admins:edit"), admins.CreateAdmin)

		// Vehicles & Rides
		vehicleAdmin := adminhandlers.NewVehiclesHandler(db, log)
		rideAdmin := adminhandlers.NewRidesHandler(db, log)

		panelAuth.GET("/vehicles", adminmiddleware.ACLMiddleware(db, "admin:vehicles:view"), vehicleAdmin.Index)
		panelAuth.POST("/vehicles/:id/verify", adminmiddleware.ACLMiddleware(db, "admin:vehicles:verify"), vehicleAdmin.Verify)
		panelAuth.GET("/rides", adminmiddleware.ACLMiddleware(db, "admin:rides:view"), rideAdmin.Index)
		panelAuth.POST("/rides/:id/cancel", adminmiddleware.ACLMiddleware(db, "admin:rides:cancel"), rideAdmin.Cancel)

		// Earnings & Withdrawals
		earningsAdmin := adminhandlers.NewEarningsHandler(db, log)
		panelAuth.GET("/earnings", adminmiddleware.ACLMiddleware(db, "admin:earnings:view"), earningsAdmin.Index)
		panelAuth.GET("/withdrawals", adminmiddleware.ACLMiddleware(db, "admin:earnings:view"), earningsAdmin.Withdrawals)
		panelAuth.POST("/withdrawals/:id/approve", adminmiddleware.ACLMiddleware(db, "admin:withdrawals:approve"), earningsAdmin.ApproveWithdrawal)
		panelAuth.POST("/withdrawals/:id/reject", adminmiddleware.ACLMiddleware(db, "admin:withdrawals:reject"), earningsAdmin.RejectWithdrawal)

		// Help Center / Issues
		helpAdmin := adminhandlers.NewHelpHandler(db, log)
		panelAuth.GET("/help", adminmiddleware.ACLMiddleware(db, "admin:help:view"), helpAdmin.Index)
		panelAuth.POST("/help/:id/reply", adminmiddleware.ACLMiddleware(db, "admin:help:reply"), helpAdmin.Reply)
		panelAuth.POST("/help/:id/close", adminmiddleware.ACLMiddleware(db, "admin:help:close"), helpAdmin.Close)

		// Settings
		settingsAdmin := adminhandlers.NewSettingsHandler(db, log)
		panelAuth.GET("/settings", adminmiddleware.ACLMiddleware(db, "admin:settings:edit"), settingsAdmin.SettingsPage)
		panelAuth.POST("/settings/password", adminmiddleware.ACLMiddleware(db, "admin:settings:edit"), settingsAdmin.UpdatePassword)

	} else {
		panel.GET("/*any", func(c *gin.Context) {
			c.String(http.StatusServiceUnavailable, "Admin panel UI is not installed. Implement templates to enable it.")
		})
	}

	return r
}
