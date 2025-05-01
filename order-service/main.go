// package main

// import (
// 	"context"
// 	"fmt"
// 	"kairos/order-service/controllers"
// 	"kairos/order-service/db"
// 	"kairos/order-service/kafka"
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/gin-gonic/gin"
// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// func main() {
// 	r := gin.Default()

// 	// Fetch environment variables
// 	dbURL := os.Getenv("DATABASE_URL")
// 	if dbURL == "" {
// 		log.Fatal("DATABASE_URL environment variable is not set")
// 	}

// 	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
// 	if kafkaBrokers == "" {
// 		log.Fatal("KAFKA_BROKERS environment variable is not set")
// 	}

// 	serverPort := os.Getenv("SERVER_PORT")
// 	if serverPort == "" {
// 		serverPort = "8008" // Default port
// 	}

// 	// Initialize DB connection
// 	pool, err := pgxpool.New(context.Background(), dbURL)
// 	if err != nil {
// 		log.Fatalf("Unable to create connection pool: %v", err)
// 	}
// 	defer pool.Close()

// 	err = pool.Ping(context.Background())
// 	if err != nil {
// 		log.Fatalf("Unable to connect to database: %v", err)
// 	}

// 	fmt.Println("Connected to the database!")

// 	// Initialize the Kafka producer
// 	kafkaConfig := kafka.NewKafkaConfig([]string{kafkaBrokers})
// 	kafkaProducer, err := kafka.NewKafkaProducer(kafkaConfig)
// 	if err != nil {
// 		log.Fatalf("Error initializing Kafka producer: %v", err)
// 	}
// 	defer kafka.CloseKafkaProducer(kafkaProducer)

// 	// Initialize queries object
// 	queries := db.New(pool)

// 	// Define routes
// 	r.GET("/", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "Order-service Started",
// 		})
// 	})

// 	// Order routes
// 	orderGroup := r.Group("/orders")
// 	{
// 		orderGroup.POST("/", controllers.CreateOrder(queries, kafkaProducer))            // POST /orders - Create Order
// 		orderGroup.GET("/", controllers.GetOrdersByConsumer(queries))                    // GET /orders?user_id=12345 - Get orders by consumer
// 		orderGroup.GET("/provider", controllers.GetOrdersByProvider(queries))            // GET /orders/provider?chef_id=provider123 - Get orders by provider
// 		orderGroup.DELETE("/:orderItemId", controllers.DeleteOrderItem(pool, queries))   // DELETE /orders/order001 - Delete Order
// 		orderGroup.PATCH("/:orderId/status", controllers.UpdateOrderItemStatus(queries)) // PATCH /orders/order987/status - Update order status
// 		orderGroup.GET("/:orderId/status", controllers.GetOrderItemStatus(queries))      // GET /orders/{orderId}/status - Get order status
// 	}

// 	// Run the server
// 	r.Run(":" + serverPort)
// }

package main

import (
	"context"
	"fmt"
	"kairos/order-service/controllers"
	"kairos/order-service/db"
	"kairos/order-service/kafka"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	r := gin.Default()

	// Enable CORS to allow all origins
	r.Use(cors.Default())

	// Fetch environment variables
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		log.Fatal("KAFKA_BROKERS environment variable is not set")
	}

	kafkaUsername := os.Getenv("KAFKA_USERNAME")
	if kafkaUsername == "" {
		log.Fatal("KAFKA_USERNAME environment variable is not set")
	}

	kafkaPassword := os.Getenv("KAFKA_PASSWORD")
	if kafkaPassword == "" {
		log.Fatal("KAFKA_PASSWORD environment variable is not set")
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8008" // Default port
	}

	// Initialize DB connection
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	defer pool.Close()

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	fmt.Println("Connected to the database!")
// Initialize Kafka config
kafkaConfig := kafka.NewKafkaConfig([]string{kafkaBrokers})

// Initialize Kafka producer
kafkaProducer, err := kafka.NewKafkaProducer(kafkaConfig)
if err != nil {
	log.Fatalf("Error initializing Kafka producer: %v", err)
}
defer kafka.CloseKafkaProducer(kafkaProducer)

// Initialize Kafka consumer
kafkaConsumer, err := kafka.NewKafkaConsumer(kafkaConfig)
if err != nil {
	log.Fatalf("Error initializing Kafka consumer: %v", err)
}
defer kafka.CloseKafkaConsumer(kafkaConsumer)

// Initialize DB query layer
queries := db.New(pool)

// Listen for reservation status events in a goroutine
go func() {
	err := kafka.ConsumeReservationStatus(
		kafkaConsumer,
		controllers.HandleReservationStatus(queries, kafkaProducer),
	)
	if err != nil {
		log.Fatalf("Reservation status consumer error: %v", err)
	}
}()


	// Define routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Order-service Started",
		})
	})

	// Order routes
	orderGroup := r.Group("/orders")
	{
		orderGroup.POST("/", controllers.CreateOrder(queries, kafkaProducer))            // POST /orders - Create Order
		orderGroup.GET("/", controllers.GetOrdersByConsumer(queries))                    // GET /orders?user_id=12345 - Get orders by consumer
		orderGroup.GET("/provider", controllers.GetOrdersByProvider(queries))            // GET /orders/provider?chef_id=provider123 - Get orders by provider
		orderGroup.DELETE("/:orderItemId", controllers.DeleteOrderItem(pool, queries))   // DELETE /orders/order001 - Delete Order
		orderGroup.PATCH("/:orderId/status", controllers.UpdateOrderItemStatus(queries,kafkaProducer)) // PATCH /orders/order987/status - Update order status
		orderGroup.GET("/:orderId/status", controllers.GetOrderItemStatus(queries))      // GET /orders/{orderId}/status - Get order status
	}

	// Run the server
	r.Run(":" + serverPort)
}
