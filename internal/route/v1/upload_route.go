// internal/route/v1/upload_route.go
package v1

import (
	"guide-me/internal/controller"
	"guide-me/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUploadRoutes(r *gin.RouterGroup) {
	upload := r.Group("/upload")
	upload.Use(middleware.AuthMiddleware())
	{
		upload.POST("/image", controller.UploadImage)
	}
}