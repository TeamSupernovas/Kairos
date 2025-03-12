package main

import (
	"context"
	"geodishdiscoveryservice/config"
	"geodishdiscoveryservice/internal/handler"
	"geodishdiscoveryservice/internal/infrastructure"
	"geodishdiscoveryservice/internal/repository"
	"geodishdiscoveryservice/internal/service"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.GetGeoDishDiscoveryServiceConfig()

	// Initialize infrastructure
	infra, err := infrastructure.NewInfrastructure(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize infrastructure: %v", err)
	}
	defer func() {
		// Close MongoDB connection
		if err := infra.DB().Disconnect(context.Background()); err != nil {
			log.Printf("Failed to disconnect MongoDB: %v", err)
		}
		// Close Kafka resources
		infra.KafkaResources().ReaderDishCreated.Close()
		infra.KafkaResources().ReaderDishUpdated.Close()
		infra.KafkaResources().ReaderDishDeleted.Close()
	}()

	// Initialize repository
	dishRepo := repository.NewDishRepository(infra.DB(), cfg.DatabaseConfig().DBName(), "dishes")

	// Initialize Kafka service
	kafkaService := service.NewKafkaService(infra.KafkaResources())

	// Initialize dish service
	dishService := service.NewDishService(dishRepo, kafkaService)

	// Subscribe to Kafka topics
	if err := dishService.SubscribeToDishTopics(context.Background()); err != nil {
		log.Fatalf("Failed to subscribe to Kafka topics: %v", err)
	}

	// Initialize handler
	dishHandler := handler.NewDishHandler(dishService)

	// Set up Gin router
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	  }))

	// Register REST API endpoint
	router.GET("/dishes/search", dishHandler.GetNearbyDishes)

	// Start the HTTP server
	log.Printf("Starting server on port %s", cfg.AppConfig().Port())
	if err := router.Run(cfg.AppConfig().Port()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
