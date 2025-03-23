package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pm-cloudify/http-server/internal/api/v1/models"
	"github.com/pm-cloudify/http-server/internal/api/v1/services"
)

func Login(c *gin.Context) {
	var loginRequest models.LoginRequest

	// getting required data from request
	err := c.ShouldBind(&loginRequest)
	if err != nil || loginRequest.Username == "" || loginRequest.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// getting token from logic
	token, err := services.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
