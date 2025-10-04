package routes

import (
	"github.com/gin-gonic/gin"
	"zipride/internal/driverAuth/handlers"
)

func DriverRoutes(c *gin.Engine) {

	driverGroup := c.Group("/driver")
	{
		driverGroup.POST("/signup", handlers.DriverSignUp)
		driverGroup.POST("/login", handlers.DriverLogin)
		driverGroup.POST("/otp/send", handlers.SendDriverOtp)
		driverGroup.POST("/otp/verify", handlers.VerifyOtpHandler)
		driverGroup.POST("/login/google", handlers.DriverGoogleLogin)
	}
}
