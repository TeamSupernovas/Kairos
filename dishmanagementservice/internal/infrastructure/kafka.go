package infrastructure

import (
	"dishmanagementservice/config"
	"fmt"

	"github.com/segmentio/kafka-go"
)

const (
	EventDishCreated = "dish-management-service.dish.created"
	EventDishUpdated = "dish-management-service.dish.updated"
	EventDishDeleted = "dish-management-service.dish.deleted"
	EventOrderCreated = "order-service.order.created"
	EventOrderUpdated = "order-service.order.updated"
	EventOrderDeleted = "order-service.order.deleted"
)

// KafkaResources holds the producer and consumer for each topic
type KafkaResources struct {
	WriterDishCreated  *kafka.Writer
	WriterDishUpdated  *kafka.Writer
	WriterDishDeleted  *kafka.Writer
	ReaderOrderCreated *kafka.Reader
	ReaderOrderUpdated *kafka.Reader
	ReaderOrderDeleted *kafka.Reader
}

// initKafka initializes Kafka producers and consumers for each topic
func initKafka(cfg config.KafkaConfig) (*KafkaResources, error) {
	brokerAddress := fmt.Sprintf("%s:%s", cfg.Host(), cfg.Port())

	// Initialize writers for each topic
	writerDishCreated := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   cfg.TopicDishCreated(),
	})

	writerDishUpdated := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   cfg.TopicDishUpdated(),
	})

	writerDishDeleted := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   cfg.TopicDishDeleted(),
	})

	// Initialize readers for each topic
	readerOrderCreated := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		GroupID: cfg.GroupID(),
		Topic:   cfg.TopicOrderCreated(),
	})

	readerOrderUpdated := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		GroupID: cfg.GroupID(),
		Topic:   cfg.TopicOrderUpdated(),
	})

	readerOrderDeleted := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		GroupID: cfg.GroupID(),
		Topic:   cfg.TopicOrderDeleted(),
	})

	// Bundle the resources into a single struct
	kafkaResources := &KafkaResources{
		WriterDishCreated: writerDishCreated,
		WriterDishUpdated: writerDishUpdated,
		WriterDishDeleted: writerDishDeleted,
		ReaderOrderCreated: readerOrderCreated,
		ReaderOrderUpdated: readerOrderUpdated,
		ReaderOrderDeleted: readerOrderDeleted,
	}

	return kafkaResources, nil
}
