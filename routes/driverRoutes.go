package routes

import (
	driverHandlers "zipride/internal/domain/driver/handlers"
	driverServices "zipride/internal/domain/driver/services"
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
			authd.POST("/complete-profile", handlers.CompleteProfileHandler)
		}

		// Document upload (multipart/form-data)
		authd.POST("/upload-document", driverServices.UploadDocument)
		authd.GET("/document-status", driverServices.GetDocumentStatus)

		// Driver operations - requires approval
		approved := driverGroup.Group("/ops").Use(middleware.Auth("driver"), middleware.RequireApprovedDriver())
		{
			approved.POST("/availability", driverServices.SetDriverAvailability)
			approved.PATCH("/location", driverServices.UpdateDriverLocation)
			approved.GET("/status", driverServices.GetDriverStatus)
		}

		// Public endpoint for finding nearby drivers
		driverGroup.GET("/nearby", driverServices.GetNearbyDrivers)

		// Example protected area once approved
		approvedExample := driverGroup.Group("/protected").Use(middleware.Auth("driver"), middleware.RequireApprovedDriver())
		{
			approvedExample.GET("/ping", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "ok"}) })
		}
	}
}
