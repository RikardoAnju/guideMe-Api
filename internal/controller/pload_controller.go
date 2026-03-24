// internal/controller/upload_controller.go
package controller

import (
	"net/http"
	"path/filepath"

	"guide-me/internal/service"

	"github.com/gin-gonic/gin"
)

var allowedExt = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

func UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "gagal membaca file: " + err.Error(),
		})
		return
	}
	defer file.Close()

	// Validasi ekstensi
	ext := filepath.Ext(header.Filename)
	if !allowedExt[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "format file tidak didukung, gunakan jpg/jpeg/png/webp",
		})
		return
	}

	// Validasi ukuran maksimal 5MB
	if header.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ukuran file maksimal 5MB",
		})
		return
	}

	url, err := service.UploadImage(file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "gambar berhasil diupload",
		"image_url": url,
	})
}