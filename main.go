package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	// docs "github.com/bjorndonald/lasgcce/docs"

	"github.com/bjorndonald/lasgcce/constants"
	"github.com/bjorndonald/lasgcce/database"
	"github.com/bjorndonald/lasgcce/docs"
	"github.com/bjorndonald/lasgcce/internal/bootstrap"
	"github.com/bjorndonald/lasgcce/internal/handlers"
	"github.com/bjorndonald/lasgcce/internal/helpers"
	"github.com/bjorndonald/lasgcce/internal/otp"
	"github.com/bjorndonald/lasgcce/internal/routes"
	"github.com/bjorndonald/lasgcce/internal/service/streaming"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"
	_ "golang.org/x/text"
)

// @title Lasgcee
// @version 1.0
// @description Swagger API documentation for Lasgcee API
// @termsOfService http://swagger.io/terms/
// @contact.name Bjorn-Donald Bassey
// @contact.email bjorndonaldb@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host http://lasgcee.alwaysdata.net
// @BasePath /api/v1
func main() {

	g := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"

	constant := constants.New()
	_ = otp.NewOTPManager()

	ctx := context.Background()

	v := constants.New()

	// Initialize the client using your apitoolkit.io generated apikey
	// apitoolkitClient, err := apitoolkit.NewClient(ctx, apitoolkit.Config{APIKey: v.APIToolkitKey})
	// if err != nil {
	// 	// Handle the error your own way
	// 	log.Println(err)
	// }

	// Parse command-line flags
	flag.Parse()

	g.Static("/assets", "./static/public")
	g.Static("/templates", "./templates")

	// Middleware
	g.Use(gin.CustomRecovery(func(c *gin.Context, recovered any) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	g.Use(gin.Logger())

	g.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	g.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // add more origins
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	g.MaxMultipartMemory = 8 << 20

	// g.Use(apitoolkit.GinMiddleware(apitoolkitClient))

	dbConfig := database.Config{
		Host:     v.DbHost,
		Port:     v.DbPort,
		Password: v.DbPassword,
		User:     v.DbUser,
		DBName:   v.DbName,
		SSLMode:  v.SSLMode,
	}
	database.Connect(&dbConfig)

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", v.DbUser, v.DbPassword, v.DbHost, v.DbPort, v.DbName, v.SSLMode)
	database.RunManualMigration(connStr)

	// Set up Swagger documentation
	docs.SwaggerInfo.BasePath = "/api/v1"
	url := ginSwagger.URL("/swagger/doc.json")
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, url))

	g.GET("/api/v1/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	v1 := g.Group("/api/v1")

	keepRunning := true

	provider, err := streaming.NewProducer(&streaming.Config{
		Verbose:   false,
		Producers: 3,
		Topic:     []string{"signup"},
		Version:   "3.8.0",
		Brokers:   "kafka:9092",
	})

	if err != nil {
		log.Fatal(err)
	}

	dependencies := bootstrap.InitializeDependencies(database.DB)
	dependencies.EventProducer = provider

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumerClient, err := streaming.NewConsumer(&streaming.Config{
		Verbose:  true,
		Version:  "3.8.0",
		Brokers:  "kafka:9092",
		Assignor: "roundrobin",
		Oldest:   true,
		Group:    "gin",
		Ctx:      ctx,
	})

	if err != nil {
		log.Fatal(err)
	}

	eventHandler := handlers.EventHandler{
		Deps: dependencies,
	}

	if err := consumerClient.Consume("signup", eventHandler.ProcessSignup); err != nil {
		log.Fatal(err)
	}

	// ///////////////////////

	// /////////////////////////////

	routes.Routes(v1, dependencies)

	g.NoRoute(func(c *gin.Context) {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf("route not found"), http.StatusNotFound)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = constant.Port
	}

	go log.Fatal(g.Run(":" + port))

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			// consumerClient.ToggleConsumptionFlow()
		}
	}
	// provider.Clear()
	// cancel()
}
