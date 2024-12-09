package main

import (
	"context"
	"fmt"
	"kairos/order-service/controllers"
	"kairos/order-service/db"
	"log"
	"net/http"
	"kairos/order-service/kafka"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	r := gin.Default()

	// Initialize DB connection
	dsn := "postgres://admin:123456@localhost:5433/orders-db?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	defer pool.Close()

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	fmt.Println("Connected to the database!")

	// Initialize the Kafka producer
	kafkaConfig := kafka.NewKafkaConfig([]string{"localhost:9092"}) 
	kafkaProducer, err := kafka.NewKafkaProducer(kafkaConfig)
	if err != nil {
		log.Fatalf("Error initializing Kafka producer: %v", err)
	}
	defer kafka.CloseKafkaProducer(kafkaProducer)

	// Initialize queries object
	queries := db.New(pool)

	// Define routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Order-service Started",
		})
	})

	// Order routes
	orderGroup := r.Group("/orders")
	{
		orderGroup.POST("/", controllers.CreateOrder(queries,kafkaProducer))                // POST /orders - Create Order
		orderGroup.GET("/", controllers.GetOrdersByConsumer(queries))         // GET /orders?consumerId=12345 - Get orders by consumer
		orderGroup.GET("/provider", controllers.GetOrdersByProvider(queries)) // GET /orders/provider?providerId=provider123 - Get orders by provider
		orderGroup.DELETE("/:orderItemId", controllers.DeleteOrderItem(pool,queries))             // DELETE /orders/order001 - Delete Order
		orderGroup.PATCH("/:orderId/status", controllers.UpdateOrderItemStatus(queries)) // PATCH /orders/order987/status - Update order status
		orderGroup.GET("/:orderId/status", controllers.GetOrderItemStatus(queries))      // GET /orders/{orderId}/status - Get order status
	}

	// Run the server
	r.Run(":8008")
}
