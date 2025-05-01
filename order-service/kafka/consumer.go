package kafka

import (
	"os"

	"github.com/IBM/sarama"
)

// Helper to load topic names from environment with fallback
func getTopicFromEnv(envVar, fallback string) string {
	if val := os.Getenv(envVar); val != "" {
		return val
	}
	return fallback
}

// ConsumeOrderPlaced consumes the "Order Placed" event
func ConsumeOrderPlaced(consumer sarama.Consumer, eventHandler func(key, value []byte)) error {
	topic := getTopicFromEnv("ORDER_PLACED_TOPIC", "order-service.order.placed")
	return ConsumeMessages(consumer, topic, eventHandler)
}

// ConsumeOrderUpdated consumes the "Order Updated" event
func ConsumeOrderUpdated(consumer sarama.Consumer, eventHandler func(key, value []byte)) error {
	topic := getTopicFromEnv("ORDER_UPDATED_TOPIC", "order-service.order.updated")
	return ConsumeMessages(consumer, topic, eventHandler)
}

// ConsumeOrderDeleted consumes the "Order Deleted" event
func ConsumeOrderDeleted(consumer sarama.Consumer, eventHandler func(key, value []byte)) error {
	topic := getTopicFromEnv("ORDER_DELETED_TOPIC", "order-service.order.deleted")
	return ConsumeMessages(consumer, topic, eventHandler)
}

// ConsumeReservationStatus consumes dish reservation status events
func ConsumeReservationStatus(consumer sarama.Consumer, eventHandler func(key, value []byte)) error {
	topic := getTopicFromEnv("RESERVATION_STATUS_TOPIC", "dish-management-service.dish.reservation-status")
	return ConsumeMessages(consumer, topic, eventHandler)
}
