package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Role sudah di-set oleh AuthMiddleware dari JWT claims
		role, exists := c.Get("role")
		if !exists || role.(string) != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"message": "Akses ditolak, hanya admin yang diizinkan"})
			c.Abort()
			return
		}
		c.Next()
	}
}