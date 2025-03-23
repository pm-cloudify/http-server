package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/pm-cloudify/http-server/internal/api/v1/handler"
)

func SetupRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/login", handler.Login)
	}
}
