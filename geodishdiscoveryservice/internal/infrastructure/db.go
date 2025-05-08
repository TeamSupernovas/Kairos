package infrastructure

import (
	"context"
	"fmt"
	"geodishdiscoveryservice/config"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func initDB(dbConfig config.DatabaseConfig) (*mongo.Client, error) {
    connectionUri := dbConfig.ConnectionURI()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
    // Set client options
    clientOptions := options.Client().ApplyURI(connectionUri).SetServerAPIOptions(serverAPI)

    // Create a new client and connect to the server
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(clientOptions)
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
