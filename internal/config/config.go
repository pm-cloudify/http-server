package config

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type AppEnv struct {
	AC_SecretKey   string
	AC_AccessKey   string
	AC_S3Bucket    string
	AC_S3_Endpoint string
	AC_S3_Region   string
}

var LoadedEnv AppEnv

func MustLoadENV() {
	if err := godotenv.Load("configs/.env"); err != nil {
		log.Panicln(err.Error())
	}
	LoadedEnv.AC_S3Bucket = os.Getenv("S3_BUCKET")
	LoadedEnv.AC_SecretKey = os.Getenv("SECRET_KEY")
	LoadedEnv.AC_AccessKey = os.Getenv("ACCESS_KEY")
	LoadedEnv.AC_S3_Endpoint = os.Getenv("S3_ENDPOINT")
	LoadedEnv.AC_S3_Region = os.Getenv("S3_REGION")
	log.Println(LoadedEnv)
}

func ConfigMiddlewares(router *gin.Engine) {
	// router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// TODO: check development and production mode and apply different policies
	// config CORS policy
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// Allow requests with no origin (like mobile apps or curl requests)
			if origin == "" {
				return true
			}

			// Allow same origin requests
			if strings.HasPrefix(origin, "http://localhost:") ||
				strings.HasPrefix(origin, "https://localhost:") {
				return true
			}

			return false
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           24 * time.Hour,
	}))
}

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
}
