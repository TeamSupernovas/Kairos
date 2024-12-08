package infrastructure

import (
	"context"
	"fmt"
	"geodishdiscoveryservice/config"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initDB(dbConfig config.DatabaseConfig) (*mongo.Client, error) {
    var uri string
    if dbConfig.Username() != "" && dbConfig.Password() != "" {
        // Use authentication if credentials are provided
        uri = fmt.Sprintf("mongodb://%s:%s@%s:%s", dbConfig.Username(), dbConfig.Password(), dbConfig.Host(), dbConfig.Port())
    } else {
        // No authentication
        uri = fmt.Sprintf("mongodb://%s:%s", dbConfig.Host(), dbConfig.Port())
    }

    // Set client options
    clientOptions := options.Client().ApplyURI(uri)

    // Create a new client and connect to the server
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
    }

    // Verify connection with Ping
    if err := client.Ping(ctx, nil); err != nil {
        return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
    }

    log.Println("Successfully connected to MongoDB")
    return client, nil
}
