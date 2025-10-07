package routes

import (
	"zipride/internal/constants"
	"zipride/internal/domain/user/services"
	"zipride/internal/middleware"
	"zipride/internal/user_Auth/handlers"

	"github.com/gin-gonic/gin"
)

// user routes
func UserRoutes(c *gin.Engine) {
	user := c.Group("/user")

	//user Forget password
	user.POST("/forget-password", handlers.ForgetPassword)
	user.POST("/veryfy-otp", handlers.VerifyForgotOTP)
	user.POST("/reset-password", handlers.ResetPassword)
	user.Use(middleware.Auth(constants.RoleUser))

	//user profile set up
	user.GET("/profile", services.GetUserProfile)
	user.PUT("/update", services.UpdateUserProfile)
	user.DELETE("/delete", services.DeleteUserProfile)
	user.POST("/logout", handlers.UserLogout)

}
