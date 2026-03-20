package controller

import (
	"net/http"

	"guide-me/internal/models"
	"guide-me/internal/service"

	"github.com/gin-gonic/gin"
)

// GET /admin/users?page=1&limit=10&search=john&role=user
func GetAllUsers(c *gin.Context) {
	var query models.GetUsersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result, err := service.GetAllUsers(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    result,
	})
}

// GET /admin/users/:id
func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"user":    user,
	})
}