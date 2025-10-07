package routes

import (
	"zipride/internal/constants"
	driverAdminHandlers "zipride/internal/domain/driverAdmin/handlers"
	driverAdminServices "zipride/internal/domain/driverAdmin/services"
	"zipride/internal/middleware"

	"github.com/gin-gonic/gin"
)

func DriverAdminRoutes(c *gin.Engine) {
	// Public admin auth
	c.POST("/admin/login", driverAdminHandlers.DriverAdminLogin)

	admin := c.Group("/admin")
	admin.Use(middleware.Auth(constants.RoleDriverAdmin))

	drivers := admin.Group("/drivers")
	{
		drivers.GET("", middleware.RequireDriverAdminPerm(constants.PermissionDriverView), driverAdminServices.GetDriversList)
		drivers.GET("/stats", middleware.RequireDriverAdminPerm(constants.PermissionDriverView), driverAdminServices.GetDriverStats)
		drivers.GET(":id", middleware.RequireDriverAdminPerm(constants.PermissionDriverView), driverAdminServices.GetDriverDetails)
		drivers.GET(":id/docs", middleware.RequireDriverAdminPerm(constants.PermissionDriverView), driverAdminServices.GetDriverDocs)
		drivers.PUT(":id/approve", middleware.RequireDriverAdminPerm(constants.PermissionDriverApprove), driverAdminServices.ApproveDriver)
		drivers.PUT(":id/reject", middleware.RequireDriverAdminPerm(constants.PermissionDriverReject), driverAdminServices.RejectDriver)
		drivers.PUT(":id/suspend", middleware.RequireDriverAdminPerm(constants.PermissionDriverSuspend), driverAdminServices.SuspendDriver)
		drivers.PUT(":id/unsuspend", middleware.RequireDriverAdminPerm(constants.PermissionDriverSuspend), driverAdminServices.UnsuspendDriver)
	}

	// staff management
	staff := admin.Group("/staff")
	{
		staff.GET("", middleware.RequireDriverAdminPerm(constants.PermissionViewStaffs), driverAdminServices.GetDriverAdminStaffList)
		staff.POST("", middleware.RequireDriverAdminPerm(constants.PermissionAddStaff), driverAdminServices.CreateDriverAdminStaff)
		staff.GET(":id", middleware.RequireDriverAdminPerm(constants.PermissionViewStaffs), driverAdminServices.GetDriverAdmin)
		staff.PUT(":id", middleware.RequireDriverAdminPerm(constants.PermissionEditStaff), driverAdminServices.UpdateDriverAdminStaff)
		staff.DELETE(":id", middleware.RequireDriverAdminPerm(constants.PermissionEditStaff), driverAdminServices.DeleteDriverAdminStaff)
		staff.POST(":id/change-password", middleware.RequireDriverAdminPerm(constants.PermissionEditStaff), driverAdminServices.ChangeDriverAdminPassword)
		staff.GET(":id/permissions", middleware.RequireDriverAdminPerm(constants.PermissionViewStaffs), driverAdminServices.GetAdminPermissions)
		staff.POST(":id/assign-role", middleware.RequireDriverAdminPerm(constants.PermissionEditStaff), driverAdminServices.AssignRoleToAdmin)
		staff.POST(":id/roles", middleware.RequireDriverAdminPerm(constants.PermissionEditStaff), driverAdminServices.AssignRoles)
	}

	// roles and permissions
	roles := admin.Group("/roles")
	{
		roles.GET("", middleware.RequireDriverAdminPerm(constants.PermissionSystemSettings), driverAdminServices.ListRoles)
		roles.POST("", middleware.RequireDriverAdminPerm(constants.PermissionSystemSettings), driverAdminServices.CreateRole)
		roles.PUT(":id", middleware.RequireDriverAdminPerm(constants.PermissionSystemSettings), driverAdminServices.UpdateRole)
		roles.DELETE(":id", middleware.RequireDriverAdminPerm(constants.PermissionSystemSettings), driverAdminServices.DeleteRole)
		roles.POST(":id/permissions", middleware.RequireDriverAdminPerm(constants.PermissionSystemSettings), driverAdminServices.SetRolePermissions)
	}

	// permissions management
	permissions := admin.Group("/permissions")
	{
		permissions.GET("", middleware.RequireDriverAdminPerm(constants.PermissionSystemSettings), driverAdminServices.GetPermissionsList)
	}

	// analytics and reporting
	analytics := admin.Group("/analytics")
	{
		analytics.GET("/drivers", middleware.RequireDriverAdminPerm(constants.PermissionDriverView), driverAdminServices.GetDriverAnalytics)
		analytics.GET("/performance", middleware.RequireDriverAdminPerm(constants.PermissionDriverView), driverAdminServices.GetDriverPerformanceMetrics)
		analytics.GET("/documents", middleware.RequireDriverAdminPerm(constants.PermissionDriverView), driverAdminServices.GetDocumentAnalytics)
		analytics.GET("/health", middleware.RequireDriverAdminPerm(constants.PermissionSystemSettings), driverAdminServices.GetSystemHealthMetrics)
	}
}
