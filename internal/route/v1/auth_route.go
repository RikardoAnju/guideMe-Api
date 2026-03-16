package v1

import (
	"guide-me/internal/controller"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", controller.Register)
		auth.POST("/login", controller.Login)
		auth.POST("/reset-password", controller.ResetPassword)
		auth.POST("/change-password", controller.ChangePassword)
	}
}