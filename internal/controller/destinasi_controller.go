// internal/controller/destinasi_controller.go
package controller

import (
	"net/http"

	"guide-me/internal/models"
	"guide-me/internal/service"

	"github.com/gin-gonic/gin"
)

func GetAllDestinasi(c *gin.Context) {
	data, err := service.GetAllDestinasi()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil semua destinasi",
		"data":    data,
	})
}

func GetDestinasiByID(c *gin.Context) {
	id := c.Param("id")

	data, err := service.GetDestinasiByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil destinasi",
		"data":    data,
	})
}

func CreateDestinasi(c *gin.Context) {
	var req models.CreateDestinasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil user_id dari context (di-set oleh middleware auth)
	userID, _ := c.Get("user_id")

	data, err := service.CreateDestinasi(req, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Destinasi berhasil dibuat",
		"data":    data,
	})
}

func UpdateDestinasi(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateDestinasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := service.UpdateDestinasi(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Destinasi berhasil diupdate",
		"data":    data,
	})
}

func DeleteDestinasi(c *gin.Context) {
	id := c.Param("id")

	if err := service.DeleteDestinasi(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Destinasi berhasil dihapus",
	})
}

func SubmitRating(c *gin.Context) {
	id := c.Param("id")

	var req models.RatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := service.SubmitRating(id, req.Rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Rating berhasil dikirim",
		"data":    data,
	})
}