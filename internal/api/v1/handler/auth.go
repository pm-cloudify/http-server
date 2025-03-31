package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pm-cloudify/http-server/internal/api/v1/models"
	"github.com/pm-cloudify/http-server/internal/api/v1/services"
)

func checkLoginRequestBind(c *gin.Context, data *models.LoginRequest) error {
	err := c.ShouldBind(data)
	if err != nil || data.Username == "" || data.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return errors.New("invalid request body")
	}
	return nil
}

func Login(c *gin.Context) {
	var loginRequest models.LoginRequest

	if err := checkLoginRequestBind(c, &loginRequest); err != nil {
		return
	}

	// getting token from logic
	token, err := services.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		if err.Error() == "database failed" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func SignIn(c *gin.Context) {
	var signInRequest models.LoginRequest

	if err := checkLoginRequestBind(c, &signInRequest); err != nil {
		return
	}

	err := services.SingIn(signInRequest.Username, signInRequest.Password)

	if err != nil {
		var status = http.StatusBadRequest
		if err.Error() == "failed to create account" {
			status = http.StatusInternalServerError
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "account created successfully"})
}
