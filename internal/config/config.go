package config

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TODO: get these information from env variables to configure server
func ConfigGinServer(router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:           ":5000",
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func ConfigGinLogger(router *gin.Engine) {
	router.Use(gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		err_msg := ""
		if len(params.ErrorMessage) > 0 {
			err_msg = fmt.Sprintf("| error: %s", params.ErrorMessage)
		}

		return fmt.Sprintf("%s | %s | %d | %s %s\n",
			params.ClientIP,
			params.Method,
			params.StatusCode,
			params.Path,
			err_msg,
		)
	}))

	router.Use(gin.Recovery())
}
