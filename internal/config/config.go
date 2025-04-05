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
	"github.com/spf13/viper"
)

type Configs struct {
	// APP configs
	Mode   string
	Secret string

	// Web server configs
	GIN_Port string

	// RMQ configs
	RMQ_Addr   string
	RMQ_User   string
	RMQ_Pass   string
	RMQ_Q_Name string

	// DB configs
	DB_User    string
	DB_Pass    string
	DB_Host    string
	DB_Name    string
	DB_SSLMode string

	// ArvanCloud S3 configs
	AC_SecretKey   string
	AC_AccessKey   string
	AC_S3Bucket    string
	AC_S3_Endpoint string
	AC_S3_Region   string
}

var AppConfigs Configs

// load app configurations
func LoadConfigs() {
	if os.Getenv("APP_ENV") != "" {
		godotenv.Load("./configs/.env." + os.Getenv("APP_ENV"))
	} else {
		godotenv.Load("./configs/.env.development")
	}

	viper.AutomaticEnv()

	// default configs if not given
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_SECRET", "your-secret") // TODO: generate a random hash for this in each run time
	viper.SetDefault("WS_PORT", "5050")

	// app
	AppConfigs.Mode = viper.GetString("APP_ENV")
	AppConfigs.Secret = viper.GetString("APP_SECRET")

	// web server config
	AppConfigs.GIN_Port = viper.GetString("WS_PORT")

	// rabbitmq configs
	AppConfigs.RMQ_Addr = viper.GetString("RMQ_ADDR")
	AppConfigs.RMQ_User = viper.GetString("RMQ_USER")
	AppConfigs.RMQ_Pass = viper.GetString("RMQ_PASS")
	AppConfigs.RMQ_Q_Name = viper.GetString("RMQ_Q_NAME")

	// db configs
	AppConfigs.DB_User = viper.GetString("DB_USER")
	AppConfigs.DB_Pass = viper.GetString("DB_PASS")
	AppConfigs.DB_Host = viper.GetString("DB_HOST")
	AppConfigs.DB_Name = viper.GetString("DB_NAME")
	AppConfigs.DB_SSLMode = viper.GetString("DB_SSL_MODE")

	// s3 configs
	AppConfigs.AC_S3Bucket = viper.GetString("S3_BUCKET")
	AppConfigs.AC_SecretKey = viper.GetString("SECRET_KEY")
	AppConfigs.AC_AccessKey = viper.GetString("ACCESS_KEY")
	AppConfigs.AC_S3_Endpoint = viper.GetString("S3_ENDPOINT")
	AppConfigs.AC_S3_Region = viper.GetString("S3_REGION")

	if AppConfigs.Mode == "development" {
		fmt.Println(AppConfigs)
	}
}

// configured application message broker
var App_MB *mb.MessageBroker

// config message broker
func MustConnectToMessageBroker() *mb.MessageBroker {
	var err error
	App_MB, err = mb.InitMessageBroker(
		AppConfigs.RMQ_Addr,
		AppConfigs.RMQ_User,
		AppConfigs.RMQ_Pass,
		AppConfigs.RMQ_Q_Name,
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
		AppConfigs.AC_AccessKey,
		AppConfigs.AC_SecretKey,
		AppConfigs.AC_S3_Region,
		AppConfigs.AC_S3_Endpoint,
	)
}

// database setup
func MustInitDatabaseConnection() {
	// TODO: use env
	err := database.InitDB(fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		AppConfigs.DB_User,
		AppConfigs.DB_Pass,
		AppConfigs.DB_Host,
		AppConfigs.DB_Name,
		AppConfigs.DB_SSLMode,
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
		Addr:           ":" + AppConfigs.GIN_Port,
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

// configures logger
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

// configure a gin and create a new engine
func ConfigAndCreateGinEngine() *gin.Engine {
	switch os.Getenv("APP_ENV") {
	case "production":
		gin.SetMode(gin.ReleaseMode)
		fmt.Println("Running in PRODUCTION mode")
	case "staging":
		gin.SetMode(gin.TestMode)
		fmt.Println("Running in STAGING mode")
	default:
		fmt.Println("Running in DEVELOPMENT mode")
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()

	return engine
}
