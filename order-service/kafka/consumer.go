package kafka

import (
	"github.com/IBM/sarama"
)

// ConsumeOrderPlaced consumes the "Order Placed" event
func ConsumeOrderPlaced(consumer sarama.Consumer, eventHandler func(key, value []byte)) error {
	topic := "order-service.order.placed"
	return ConsumeMessages(consumer, topic, eventHandler)
}

// ConsumeOrderUpdated consumes the "Order Updated" event
func ConsumeOrderUpdated(consumer sarama.Consumer, eventHandler func(key, value []byte)) error {
	topic := "order-service.order.updated"
	return ConsumeMessages(consumer, topic, eventHandler)
}

// ConsumeOrderDeleted consumes the "Order Deleted" event
func ConsumeOrderDeleted(consumer sarama.Consumer, eventHandler func(key, value []byte)) error {
	topic := "order-service.order.deleted"
	return ConsumeMessages(consumer, topic, eventHandler)
}
