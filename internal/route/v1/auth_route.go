package v1

import (
	"guide-me/internal/controller"
	"guide-me/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", controller.Register)
		auth.POST("/login", controller.Login)
		auth.GET("/verify-email", controller.VerifyEmail)
		auth.POST("/resend-verification", controller.ResendVerificationEmail)
		auth.POST("/reset-password", controller.ResetPassword)  
		auth.POST("/logout", middleware.AuthMiddleware(), controller.Logout)
		auth.POST("/verify-otp", controller.VerifyOTP)           
		auth.POST("/change-password", controller.ChangePassword) 
	}
}
