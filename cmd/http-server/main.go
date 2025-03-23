package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pm-cloudify/http-server/internal/config"
)

func main() {
	router := gin.New()

	config.ConfigGinLogger(router)
	server := config.ConfigGinServer(router)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	server.ListenAndServe()
}
