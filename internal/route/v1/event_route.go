// internal/route/v1/event_route.go
package v1

import (
	"guide-me/internal/controller"
	"guide-me/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupEventRoutes(r *gin.RouterGroup) {
	event := r.Group("/event")
	{
		// Public
		event.GET("", controller.GetAllEvent)
		event.GET("/:id", controller.GetEventByID)

		// Protected (perlu login)
		protected := event.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("", controller.CreateEvent)
			protected.PUT("/:id", controller.UpdateEvent)
			protected.DELETE("/:id", controller.DeleteEvent)
		}
	}
}