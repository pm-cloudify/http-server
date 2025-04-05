package main

import (
	v1 "github.com/pm-cloudify/http-server/internal/api/v1"
	"github.com/pm-cloudify/http-server/internal/config"
	"github.com/pm-cloudify/shared-libs/psql"
)

func main() {
	// create an engine
	router := config.ConfigAndCreateGinEngine()

	// loading env
	config.LoadConfigs()

	// config router
	config.ConfigGinLogger(router)
	config.ConfigMiddlewares(router)

	// setup routes
	v1.SetupRoutes(router)

	// Initialize a database
	config.MustInitDatabaseConnection()
	defer psql.CloseDB()

	// Initialize connection to s3
	config.InitS3Connection()

	// Initialize connection to RabbitMQ
	app_mb := config.MustConnectToMessageBroker()
	defer app_mb.Close()

	// config and run server
	server := config.ConfigGinServer(router)
	server.ListenAndServe()
}
