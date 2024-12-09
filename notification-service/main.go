package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gocql/gocql"
	"github.com/gorilla/websocket"
	"github.com/go-redis/redis/v8"
	"context"
)

var (
	// WebSocket upgrader
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// Active WebSocket connections mapped by userID
	clients   = make(map[string]*websocket.Conn) // userID -> *websocket.Conn
	clientsMu sync.Mutex

	// Redis client
	rdb *redis.Client
	// Redis context
	ctx = context.Background()

	// Cassandra session
	cassandraSession *gocql.Session
)

type Notification struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func main() {
	// Initialize Cassandra
	log.Println("Initializing Cassandra...")
	initCassandra()
	defer cassandraSession.Close()

	// Initialize Redis
	log.Println("Initializing Redis...")
	initRedis()

	// Kafka consumer setup
	log.Println("Setting up Kafka consumer...")
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
		"group.id":          "notification-service",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %s", err)
	}
	defer consumer.Close()

	log.Println("Subscribing to Kafka topics: order_events, dish_events...")
	consumer.SubscribeTopics([]string{"order_events", "dish_events"}, nil)

	// Start Kafka consumer in a goroutine
	go kafkaConsumer(consumer)

	// Start WebSocket server
	log.Println("Starting WebSocket server on :8080...")
	http.HandleFunc("/ws", handleWebSocket)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initCassandra() {
	cluster := gocql.NewCluster("cassandra")
	cluster.Keyspace = "notifications"
	cluster.Consistency = gocql.Quorum

	var session *gocql.Session
	var err error
	for retries := 40; retries > 0; retries-- {
		session, err = cluster.CreateSession()
		if err == nil {
			break
		}
		log.Printf("Retrying Cassandra connection: %v", err)
		time.Sleep(10 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}

	cassandraSession = session
	log.Println("Connected to Cassandra")
	createKeyspaceIfNotExists()
}

func initRedis() {
	// Retry parameters
	maxRetries := 50
	retryInterval := 10 * time.Second

	// Attempt to create a Redis client and ping the server
	var err error
	for retries := maxRetries; retries > 0; retries-- {
		rdb = redis.NewClient(&redis.Options{
			Addr: "redis:6379", // Use Redis container name instead of localhost
		})

		// Check Redis connection
		_, err = rdb.Ping(ctx).Result()
		if err == nil {
			log.Println("Connected to Redis")
			return
		}

		log.Printf("Failed to connect to Redis (attempt %d/%d): %v", maxRetries-retries+1, maxRetries, err)
		time.Sleep(retryInterval)
	}

	log.Fatalf("Failed to connect to Redis after %d retries: %v", maxRetries, err)
}

func createKeyspaceIfNotExists() {
	var exists bool
	err := cassandraSession.Query(
		`SELECT keyspace_name FROM system_schema.keyspaces WHERE keyspace_name = 'notifications' LIMIT 1`,
	).Scan(&exists)
	if err != nil || !exists {
		log.Println("Keyspace 'notifications' does not exist. Creating it...")
		err = cassandraSession.Query(
			`CREATE KEYSPACE IF NOT EXISTS notifications WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 3}`,
		).Exec()
		if err != nil {
			log.Fatalf("Failed to create keyspace: %v", err)
		}
		log.Println("Keyspace 'notifications' created successfully")
	} else {
		log.Println("Keyspace 'notifications' already exists")
	}
}

func kafkaConsumer(consumer *kafka.Consumer) {
	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer error: %v (%v)", err, msg)
			continue
		}

		var notification Notification
		if err := json.Unmarshal(msg.Value, &notification); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		log.Printf("Received notification: %v", notification)

		// Process notification
		processNotification(notification)
	}
}

func processNotification(notification Notification) {
	log.Printf("Processing notification for user %s", notification.UserID)

	// Check if the user is online via Redis
	clientID, err := rdb.Get(ctx, notification.UserID).Result()
	if err != nil {
		log.Printf("User %s is offline, sending push notification...", notification.UserID)
		sendPushNotification(notification)
		return
	}
	log.Printf("%s",clientID)
	// If user is online, send WebSocket message
	conn := getConnectionByID(notification.UserID)
	if conn != nil {
		if err := conn.WriteJSON(notification); err != nil {
			log.Printf("Failed to send WebSocket notification: %v", err)
			deleteConnection(conn)
		} else {
			log.Printf("Notification sent to user %s via WebSocket", notification.UserID)
		}
	} if conn==nil {
		log.Printf("No WebSocket connection found for user %s", notification.UserID)
	}
}

func sendPushNotification(notification Notification) {
	log.Printf("Sending push notification to %s: %s", notification.UserID, notification.Message)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket: %v", err)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		log.Println("Missing user_id in WebSocket connection")
		conn.Close()
		return
	}

	// Store the connection in Redis with expiration (e.g., 30 minutes)
	err = rdb.SetEX(ctx, userID, conn.RemoteAddr().String(), 30*time.Minute).Err()
	if err != nil {
		log.Printf("Failed to store WebSocket connection in Redis: %v", err)
		conn.Close()
		return
	}

	clientsMu.Lock()
	clients[userID] = conn // store the connection by userID
	clientsMu.Unlock()

	log.Printf("User %s connected via WebSocket",userID)

	// Handle WebSocket disconnection
	defer func() {
		clientsMu.Lock()
		delete(clients, userID) // delete the connection by userID
		clientsMu.Unlock()
		conn.Close()
		log.Printf("User %s disconnected from WebSocket", userID)

		// Remove the connection from Redis
		err = rdb.Del(ctx, userID).Err()
		if err != nil {
			log.Printf("Failed to remove WebSocket connection from Redis: %v", err)
		}
	}()

	// Listen for messages (in this case, handle heartbeats or other types of messages)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
	}
}

func getConnectionByID(userID string) *websocket.Conn {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	return clients[userID]
}

func deleteConnection(conn *websocket.Conn) {
	conn.Close()
}
