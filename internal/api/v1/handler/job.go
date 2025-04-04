package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pm-cloudify/http-server/internal/api/v1/models"
	"github.com/pm-cloudify/http-server/internal/api/v1/services"
)

// run uploaded file
func RequestRun(c *gin.Context) {
	var body models.JobRequest
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	if err := services.SendRunRequest(username, body.FileID); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": err.Error(),
		})
		return
	}
}
