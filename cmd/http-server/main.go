package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/pm-cloudify/http-server/internal/api/v1"
	"github.com/pm-cloudify/http-server/internal/config"
	"github.com/pm-cloudify/http-server/internal/database"
	"github.com/pm-cloudify/http-server/pkg/acs3"
)

func main() {
	router := gin.New()

	// loading env
	config.MustLoadENV()

	// config router
	config.ConfigGinLogger(router)
	config.ConfigMiddlewares(router)

	// ping server
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// setup routes
	v1.SetupRoutes(router)

	// Initialize a database
	// TODO: use env
	err := database.InitDB("postgres://test_user:Sample1234Pass@localhost:5432/pm_cloudify_db?sslmode=disable")
	if err != nil {
		log.Fatal("database connection failed!")
	}
	defer database.CloseDB()

	// Initialize connection to s3
	acs3.InitConnection()
	acs3.GetBuckets()

	// config and run server
	server := config.ConfigGinServer(router)
	server.ListenAndServe()
}
