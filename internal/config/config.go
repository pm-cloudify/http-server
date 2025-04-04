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
	"github.com/pm-cloudify/shared-libs/acs3"
	"github.com/pm-cloudify/shared-libs/mb"
	database "github.com/pm-cloudify/shared-libs/psql"
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

// configured application message broker
var App_MB *mb.MessageBroker

// config message broker
func MustConnectToMessageBroker() *mb.MessageBroker {
	var err error
	App_MB, err = mb.InitMessageBroker(
		os.Getenv("RMQ_ADDR"),
		os.Getenv("RMQ_USER"),
		os.Getenv("RMQ_PASS"),
		os.Getenv("RMQ_Q_NAME"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	log.Println("RabbitMQ connection established")

	return App_MB
}

// connect to s3
func InitS3Connection() {
	// Initialize connection to s3
	acs3.InitConnection(
		LoadedEnv.AC_AccessKey,
		LoadedEnv.AC_SecretKey,
		LoadedEnv.AC_S3_Region,
		LoadedEnv.AC_S3_Endpoint,
	)
}

// database setup
func MustInitDatabaseConnection() {
	// TODO: use env
	err := database.InitDB(fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	))
	if err != nil {
		log.Fatal("database connection failed!")
	}
}

// config middlewares
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
	// ping server
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

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
