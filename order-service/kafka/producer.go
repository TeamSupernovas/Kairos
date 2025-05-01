package kafka

import (
	"os"

	"github.com/IBM/sarama"
)

var (
	OrderPlacedTopic     = getEnv("ORDER_PLACED_TOPIC", "order-service.order.placed")
	NotificationTopic    = getEnv("NOTIFICATION_TOPIC", "notification-service")
)

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// OrderPlacedEvent is the structure for the order placed event
type OrderPlacedEvent struct {
	OrderID  string `json:"order_id"`
	DishID   string `json:"dish_id"`
	Portions int32  `json:"portions"`
}

// PublishOrderPlaced publishes the "Order Placed" event
func PublishOrderPlaced(producer sarama.SyncProducer, orderPlacedEvent OrderPlacedEvent) error {
	key := orderPlacedEvent.OrderID
	return PublishMessage(producer, OrderPlacedTopic, key, orderPlacedEvent)
}

// NotificationEvent defines the payload for user notifications
type NotificationEvent struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func SendNotification(producer sarama.SyncProducer, notif NotificationEvent) error {
	return PublishMessage(producer, NotificationTopic, notif.UserID, notif)
}
