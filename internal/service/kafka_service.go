package service

import (
	"context"
	"dishmanagementservice/internal/infrastructure"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaService struct {
	kafkaResources *infrastructure.KafkaResources
}

// NewKafkaService initializes a new KafkaService
func NewKafkaService(kafkaResources *infrastructure.KafkaResources) *KafkaService {
	return &KafkaService{
		kafkaResources: kafkaResources,
	}
}

func (s *KafkaService) PublishEvent(eventType string, eventData interface{}) error {
	message, err := json.Marshal(eventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %v", err)
	}

	var writer *kafka.Writer
	switch eventType {
	case infrastructure.EventDishCreated:
		writer = s.kafkaResources.WriterDishCreated
	case infrastructure.EventDishUpdated:
		writer = s.kafkaResources.WriterDishUpdated
	case infrastructure.EventDishDeleted:
		writer = s.kafkaResources.WriterDishDeleted
	default:
		return fmt.Errorf("unknown event type: %s", eventType)
	}

	return writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(eventType),
		Value: message,
	})
}


func (s *KafkaService) SubscribeToTopic(topic string, handler func(kafka.Message) error) error {
	var reader *kafka.Reader
	switch topic {
	case infrastructure.EventOrderCreated:
		reader = s.kafkaResources.ReaderOrderCreated
	case infrastructure.EventOrderUpdated:
		reader = s.kafkaResources.ReaderOrderUpdated
	case infrastructure.EventOrderDeleted:
		reader = s.kafkaResources.ReaderOrderDeleted
	default:
		return fmt.Errorf("unknown topic: %s", topic)
	}

	// Consume messages from the topic
	go func() {
		for {
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				fmt.Printf("Error reading message from %s: %v\n", topic, err)
				continue
			}

			// Process the message using the handler function
			if handlerErr := handler(msg); handlerErr != nil {
				fmt.Printf("Error handling message from %s: %v\n", topic, handlerErr)
			}
		}
	}()
	
	return nil
}
