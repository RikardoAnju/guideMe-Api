package v1

import (
	"github.com/gin-gonic/gin"
	"guide-me/internal/controller"
	"guide-me/internal/middleware"
)

func SetupAdminRoutes(r *gin.RouterGroup) {
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	{
		admin.GET("/users",      controller.GetAllUsers)
		admin.GET("/users/:id",  controller.GetUserByID)
	}
}