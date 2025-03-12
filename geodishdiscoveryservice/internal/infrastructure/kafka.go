package infrastructure

import (
	"crypto/tls"
	"fmt"
	"geodishdiscoveryservice/config"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

const (
	EventDishCreated = "dish-management-service.dish.created"
	EventDishUpdated = "dish-management-service.dish.updated"
	EventDishDeleted = "dish-management-service.dish.deleted"
)

// KafkaResources holds the producer and consumer for each topic
type KafkaResources struct {
	ReaderDishCreated *kafka.Reader
	ReaderDishUpdated *kafka.Reader
	ReaderDishDeleted *kafka.Reader
}

// initKafka initializes Kafka producers and consumers for each topic
func initKafka(cfg config.KafkaConfig) (*KafkaResources, error) {
	brokerAddress := fmt.Sprintf("%s:%s", cfg.Host(), cfg.Port())

	mechanism := plain.Mechanism {Username: cfg.Username(), Password: cfg.Password()}
	
	dialer := &kafka.Dialer{
		Timeout:       60 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
		TLS: 			&tls.Config{},
	}


	// Initialize readers for each topic
	readerDishCreated := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		GroupID: cfg.GroupID(),
		Topic:   cfg.TopicDishCreated(),
		Dialer:  dialer,
	})

	readerDishUpdated := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		GroupID: cfg.GroupID(),
		Topic:   cfg.TopicDishUpdated(),
	})

	readerDishDeleted := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		GroupID: cfg.GroupID(),
		Topic:   cfg.TopicDishDeleted(),
	})

	// Bundle the resources into a single struct
	kafkaResources := &KafkaResources{
		ReaderDishCreated: readerDishCreated,
		ReaderDishUpdated: readerDishUpdated,
		ReaderDishDeleted: readerDishDeleted,
	}

	return kafkaResources, nil
}
