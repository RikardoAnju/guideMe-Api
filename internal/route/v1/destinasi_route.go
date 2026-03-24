// internal/route/v1/destinasi_route.go
package v1

import (
	"guide-me/internal/controller"
	"guide-me/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupDestinasiRoutes(r *gin.RouterGroup) {
	destinasi := r.Group("/destinasi")
	{
		// Public
		destinasi.GET("", controller.GetAllDestinasi)
		destinasi.GET("/:id", controller.GetDestinasiByID)
		destinasi.POST("/:id/rating", controller.SubmitRating)

		// Protected (perlu login)
		protected := destinasi.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("", controller.CreateDestinasi)
			protected.PUT("/:id", controller.UpdateDestinasi)
			protected.DELETE("/:id", controller.DeleteDestinasi)
		}
	}
}
