package v1

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	v1 := r.Group("/api/v1")
	{
		SetupAuthRoutes(v1)
		SetupUserRoutes(v1)
	}
}