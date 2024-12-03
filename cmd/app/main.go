package main

import (
	"dishmanagementservice/config"
	"dishmanagementservice/internal/handler"
	"dishmanagementservice/internal/infrastructure"
	"dishmanagementservice/internal/repository"
	"dishmanagementservice/internal/service"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	// Load configurations
	cfg := config.GetDishManagementServiceConfig()

	// Initialize infrastructure
	infra, err := infrastructure.NewInfrastructure(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize infrastructure: %v", err)
	}
	defer infra.DB().Close()
	// Close all Kafka producers (writers)
	defer infra.KafkaResources().WriterDishCreated.Close()
	defer infra.KafkaResources().WriterDishUpdated.Close()
	defer infra.KafkaResources().WriterDishDeleted.Close()

	// Close all Kafka consumers (readers)
	defer infra.KafkaResources().ReaderOrderCreated.Close()
	defer infra.KafkaResources().ReaderOrderUpdated.Close()
	defer infra.KafkaResources().ReaderOrderDeleted.Close()

	// Initialize repository, service, and handler
	geoService := service.NewGeoLocationService(infra.AWSConfig())
	kafkaService := service.NewKafkaService(infra.KafkaResources())
	dishRepo := repository.NewPostgresDishRepository(infra.DB())
	dishService := service.NewDishService(dishRepo, geoService, kafkaService)
	dishHandler := handler.NewDishHandler(dishService)

	// Set up Gin router
	router := gin.Default()

	// Register the dish handler endpoint
	router.POST("/dishes", dishHandler.AddDish)
	router.PUT("/dishes/:dishId", dishHandler.UpdateDish)
	router.GET("/dishes/:dishId", dishHandler.GetDishByID)
	router.PATCH("/dishes/:dishId/delete", dishHandler.DeleteDishById)

	// Start the HTTP server
	err = router.Run(cfg.AppConfig().Port())
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}




