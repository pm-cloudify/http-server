package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/pm-cloudify/http-server/internal/api/v1"
	"github.com/pm-cloudify/http-server/internal/config"
)

func main() {
	router := gin.New()

	// config router
	config.ConfigGinLogger(router)

	// ping server
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// setup routes
	v1.SetupRoutes(router)

	// config and run server
	server := config.ConfigGinServer(router)
	server.ListenAndServe()
}
