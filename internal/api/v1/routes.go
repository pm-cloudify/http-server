package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/pm-cloudify/http-server/internal/api/common/middleware"
	"github.com/pm-cloudify/http-server/internal/api/v1/handler"
)

// setup v1 api routes
func SetupRoutes(router *gin.Engine) {

	router.POST("api/v1/login", handler.Login)

	authorized_v1 := router.Group("/api/v1")
	authorized_v1.Use(middleware.AuthMiddleware())
	{
		authorized_v1.POST("/upload", handler.Upload)
	}
}
