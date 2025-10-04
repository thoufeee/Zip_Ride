package routes

import (
	"zipride/internal/driverAuth/handlers"

	"github.com/gin-gonic/gin"
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
