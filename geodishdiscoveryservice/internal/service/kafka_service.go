package service

import (
	"context"
	"encoding/json"
	"fmt"
	"geodishdiscoveryservice/internal/infrastructure"

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


func (s *KafkaService) SubscribeToTopic(topic string, handler func(kafka.Message) error) error {
	var reader *kafka.Reader
	switch topic {
	case infrastructure.EventDishCreated:
		reader = s.kafkaResources.ReaderDishCreated
	case infrastructure.EventDishUpdated:
		reader = s.kafkaResources.ReaderDishUpdated
	case infrastructure.EventDishDeleted:
		reader = s.kafkaResources.ReaderDishDeleted
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

// ProcessDishCreated handles messages for the "Dish Created" event
func (s *KafkaService) ProcessDishCreated(msg kafka.Message) (map[string]interface{}, error) {
	fmt.Printf("Processing Dish Created event: %s\n", string(msg.Value))

	var eventData map[string]interface{}
	if err := json.Unmarshal(msg.Value, &eventData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Dish Created event: %v", err)
	}

	fmt.Printf("Dish Created event data: %+v\n", eventData)
	return eventData, nil
}

// ProcessDishUpdated handles messages for the "Dish Updated" event
func (s *KafkaService) ProcessDishUpdated(msg kafka.Message) (map[string]interface{}, error) {
	fmt.Printf("Processing Dish Updated event: %s\n", string(msg.Value))

	var eventData map[string]interface{}
	if err := json.Unmarshal(msg.Value, &eventData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Dish Updated event: %v", err)
	}

	fmt.Printf("Dish Updated event data: %+v\n", eventData)
	return eventData, nil
}

// ProcessDishDeleted handles messages for the "Dish Deleted" event
func (s *KafkaService) ProcessDishDeleted(msg kafka.Message) (map[string]interface{}, error) {
	fmt.Printf("Processing Dish Deleted event: %s\n", string(msg.Value))

	var eventData map[string]interface{}
	if err := json.Unmarshal(msg.Value, &eventData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Dish Deleted event: %v", err)
	}

	fmt.Printf("Dish Deleted event data: %+v\n", eventData)
	return eventData, nil
}
