package v1

import (
	"github.com/gin-gonic/gin"
	"guide-me/internal/controller"
	"guide-me/internal/middleware"
)

func SetupUserRoutes(r *gin.RouterGroup) {
	user := r.Group("/user")
	user.Use(middleware.AuthMiddleware())
	{
		user.GET("/profile", controller.GetProfile)
		user.PUT("/profile", controller.UpdateProfile)
	}
}
