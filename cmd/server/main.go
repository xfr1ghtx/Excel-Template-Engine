package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/stepanpotapov/Excel-Template-Engine/internal/config"
	"github.com/stepanpotapov/Excel-Template-Engine/internal/handlers"
	"github.com/stepanpotapov/Excel-Template-Engine/internal/repository"
	"github.com/stepanpotapov/Excel-Template-Engine/internal/services"
	"github.com/stepanpotapov/Excel-Template-Engine/internal/utils"
)

func main() {
	// Initialize logger
	if err := utils.InitLogger("logs.txt"); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer utils.CloseLogger()

	utils.LogInfo("Starting Acts Service...")

	// Load configuration
	cfg := config.Load()

	// Connect to MongoDB
	mongoClient := repository.ConnectMongoDB(cfg)
	defer func() {
		if err := mongoClient.Disconnect(); err != nil {
			utils.LogError("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Initialize repositories
	actRepo := repository.NewActRepository(mongoClient)

	// Initialize services
	excelService := services.NewExcelService(cfg)
	actService := services.NewActService(actRepo, excelService, cfg)

	// Initialize handlers
	actHandler := handlers.NewActHandler(actService, cfg)

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "acts-service",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		act := api.Group("/act")
		{
			act.POST("/create", actHandler.CreateAct)
			act.GET("/generate", actHandler.GenerateAct)
			act.GET("/download/:filename", actHandler.DownloadAct)
		}
	}

	// Start server in a goroutine
	go func() {
		addr := cfg.ServerHost + ":" + cfg.ServerPort
		utils.LogInfo("Server starting on %s", addr)
		if err := router.Run(addr); err != nil {
			utils.LogError("Failed to start server: %v", err)
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	utils.LogInfo("Shutting down server...")
}

