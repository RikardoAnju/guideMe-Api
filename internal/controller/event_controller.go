// internal/handler/event_handler.go
package controller

import (
	"guide-me/internal/models"
	"guide-me/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllEvent(c *gin.Context) {
	events, err := service.GetAllEvent()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "berhasil mengambil semua event",
		"data":    events,
	})
}

func GetEventByID(c *gin.Context) {
	id := c.Param("id")

	event, err := service.GetEventByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "berhasil mengambil event",
		"data":    event,
	})
}

func CreateEvent(c *gin.Context) {
	var req models.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "request tidak valid: " + err.Error(),
		})
		return
	}

	// Ambil user ID dari JWT middleware (sesuaikan key-nya dengan middleware kamu)
	createdBy, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "unauthorized",
		})
		return
	}

	event, err := service.CreateEvent(req, createdBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "event berhasil dibuat",
		"data":    event,
	})
}

func UpdateEvent(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "request tidak valid: " + err.Error(),
		})
		return
	}

	event, err := service.UpdateEvent(id, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "event tidak ditemukan" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "event berhasil diupdate",
		"data":    event,
	})
}

func DeleteEvent(c *gin.Context) {
	id := c.Param("id")

	if err := service.DeleteEvent(id); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "event tidak ditemukan" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "event berhasil dihapus",
	})
}