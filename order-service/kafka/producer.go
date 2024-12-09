package kafka

import (
	
	"github.com/IBM/sarama"
)

// OrderPlacedEvent is the structure for the order placed event
type OrderPlacedEvent struct {
	OrderID    string `json:"order_id"`
	DishID 	   string `json:"dish_id"`
	Portions int32 	`json:"portions"`
}

// PublishOrderPlaced publishes the "Order Placed" event
func PublishOrderPlaced(producer sarama.SyncProducer, orderPlacedEvent OrderPlacedEvent) error {
	topic := "order-service.order.placed"
	key := orderPlacedEvent.OrderID
	return PublishMessage(producer, topic, key, orderPlacedEvent)
}
