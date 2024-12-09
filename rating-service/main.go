package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"kairos/rating-service/controllers"
	"kairos/rating-service/db"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Database connection
	dsn := "postgres://admin:123456@localhost:5434/rating-db?sslmode=disable"

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
	r := gin.Default() // Create a Gin router

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Rating-service Started"})
	})

	// Initialize queries object
	queries := db.New(pool)

	// CRUD Routes for Ratings
	r.POST("/ratings", func(c *gin.Context) { controllers.CreateRating(c, queries) })       // Create Rating
	r.GET("/ratings/:id", func(c *gin.Context) { controllers.GetRating(c, queries) })       // Get Rating by ID
	r.GET("/ratings", func(c *gin.Context) { controllers.ListRatings(c, queries) })         // List All Ratings
	r.PUT("/ratings/:id", func(c *gin.Context) { controllers.UpdateRating(c, queries) })    // Update Rating by ID
	r.DELETE("/ratings/:id", func(c *gin.Context) { controllers.DeleteRating(c, queries) }) // Delete Rating by ID

	// Start the server on port 8080
	if err := r.Run(":8080"); err != nil {
		panic("Failed to start the server: " + err.Error())
	}
}
