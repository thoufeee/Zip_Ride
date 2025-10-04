package routes

import (
	"zipride/internal/constants"
	"zipride/internal/domain/user_Admin/services"
	"zipride/internal/middleware"

	"github.com/gin-gonic/gin"
)

// admin routes

func AdminRoutes(c *gin.Engine) {

	admin := c.Group("/admin")

	admin.Use(middleware.Auth(constants.RoleAdmin))

	admin.GET("/allusers", services.GetAllUsers)
}
