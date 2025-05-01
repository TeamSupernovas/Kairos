package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"kairos/rating-service/controllers"
	"kairos/rating-service/db"
	"kairos/rating-service/kafka"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env vars")
	}

	// Database connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("PostgreSQL connection failed: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("PostgreSQL ping failed: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	// Kafka setup
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		log.Fatal("KAFKA_BROKERS not set")
	}
	err = kafka.InitKafkaProducer(strings.Split(brokers, ","))
	if err != nil {
		log.Fatalf(" Kafka producer init failed: %v", err)
	}
	log.Println(" Kafka producer initialized")

	// Gin router setup
	r := gin.Default()

	// Enable CORS to allow all origins
	r.Use(cors.Default())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Rating-service is running"})
	})

	queries := db.New(pool)

	// Routes
	r.POST("/ratings", func(c *gin.Context) { controllers.CreateRating(c, queries) })
	r.GET("/ratings/:id", func(c *gin.Context) { controllers.GetRating(c, queries) })
	r.GET("/ratings", func(c *gin.Context) { controllers.ListRatings(c, queries) })
	r.PUT("/ratings/:id", func(c *gin.Context) { controllers.UpdateRating(c, queries) })
	r.DELETE("/ratings/:id", func(c *gin.Context) { controllers.DeleteRating(c, queries) })

	r.GET("/dishes/:dishId/ratings", func(c *gin.Context) {
		controllers.ListRatingsForDish(c, queries)
	})
	r.GET("/users/:userId/ratings", func(c *gin.Context) {
		controllers.ListRatingsByUser(c, queries)
	})
	r.GET("/ratings/chef/:chefId", func(c *gin.Context) {
		controllers.ListRatingsByChef(c, queries)
	})
	r.GET("/ratings/chef/:chefId/average", func(c *gin.Context) {
		controllers.GetChefAverageRating(c, queries)
	})
	r.GET("/ratings/dish/:dishId/average", func(c *gin.Context) {
		controllers.GetDishAverageRating(c, queries)
	})

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 Server listening on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
