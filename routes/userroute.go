package routes

import (
	"zipride/internal/constants"
	"zipride/internal/domain/auth/handlers"
	"zipride/internal/domain/user/services"
	"zipride/internal/middleware"

	"github.com/gin-gonic/gin"
)

// user routes
func UserRoutes(c *gin.Engine) {
	user := c.Group("/user")

	user.Use(middleware.Auth(constants.RoleUser))

	user.GET("/profile", services.GetUserProfile)
	user.PUT("/update", services.UpdateUserProfile)
	user.DELETE("/delete", services.DeleteUserProfile)
	user.POST("/logout", handlers.UserLogout)
}
