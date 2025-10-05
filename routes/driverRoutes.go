package routes

import (
	driverHandlers "zipride/internal/domain/driver/handlers"
	"zipride/internal/driverAuth/handlers"
	"zipride/internal/middleware"

	"github.com/gin-gonic/gin"
)

func DriverRoutes(c *gin.Engine) {

	driverGroup := c.Group("/driver")
	{
		// Phone-first OTP endpoints
		driverGroup.POST("/otp/send", handlers.SendDriverOtpHandler)
		driverGroup.POST("/otp/verify", handlers.VerifyDriverOtpHandler)

		// Google sign-in
		driverGroup.POST("/google", handlers.GoogleSignIn)

		// Token lifecycle
		driverGroup.POST("/auth/refresh", handlers.RefreshToken)
		driverGroup.POST("/auth/logout", handlers.Logout)

		// Optional email/password endpoints (secondary)
		driverGroup.POST("/signup", handlers.DriverSignUp)
		driverGroup.POST("/login", handlers.DriverLogin)

		// Onboarding - requires driver auth but not yet approved
		authd := driverGroup.Group("/me").Use(middleware.Auth("driver"))
		{
			authd.GET("", driverHandlers.Me)
			authd.POST("/profile", driverHandlers.UpdateProfile)
			authd.POST("/vehicle", driverHandlers.UpsertVehicle)
			authd.POST("/docs", driverHandlers.UploadDocs)
		}

		// Example protected area once approved
		approved := driverGroup.Group("/protected").Use(middleware.Auth("driver"), middleware.RequireApprovedDriver())
		{
			approved.GET("/ping", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "ok"}) })
		}
	}
}
