package routes

import (
	"zipride/internal/auth/handlers"
	"zipride/internal/auth/services"
	"zipride/internal/constants"
	"zipride/internal/middleware"

	"github.com/gin-gonic/gin"
)

// user routes
func UserRoutes(c *gin.Engine) {
	user := c.Group("/user")

	user.Use(middleware.Auth(constants.RoleUser))

	user.GET("/profile", services.UserProfile)
	user.POST("/logout", handlers.UserLogout)
}
