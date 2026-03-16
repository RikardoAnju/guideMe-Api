package v1

import (
	"github.com/gin-gonic/gin"
	"guide-me/internal/controller"
)

func SetupAuthRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", controller.Register)
		auth.POST("/login", controller.Login)
		auth.GET("/verify-email", controller.VerifyEmail)
		auth.POST("/resend-verification", controller.ResendVerificationEmail)
		auth.POST("/reset-password", controller.ResetPassword)
		auth.POST("/change-password", controller.ChangePassword)
	}
}
