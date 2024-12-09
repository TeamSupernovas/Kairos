package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
)

// KafkaConfig holds Kafka connection details
type KafkaConfig struct {
	Brokers []string
}

// NewKafkaConfig initializes the Kafka configuration
func NewKafkaConfig(brokers []string) *KafkaConfig {
	return &KafkaConfig{
		Brokers: brokers,
	}
}

// NewKafkaProducer initializes and returns a Kafka producer
func NewKafkaProducer(config *KafkaConfig) (sarama.SyncProducer, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Retry.Max = 5
	kafkaConfig.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(config.Brokers, kafkaConfig)
	if err != nil {
		log.Println("Error creating Kafka producer:", err)
		return nil, err
	}
	log.Println("Kafka producer created successfully")
	return producer, nil
}

// NewKafkaConsumer initializes and returns a Kafka consumer
func NewKafkaConsumer(config *KafkaConfig) (sarama.Consumer, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(config.Brokers, kafkaConfig)
	if err != nil {
		log.Println("Error creating Kafka consumer:", err)
		return nil, err
	}
	log.Println("Kafka consumer created successfully")
	return consumer, nil
}

// CloseKafkaProducer safely closes the Kafka producer
func CloseKafkaProducer(producer sarama.SyncProducer) error {
	err := producer.Close()
	if err != nil {
		log.Println("Error closing Kafka producer:", err)
		return err
	}
	log.Println("Kafka producer closed successfully")
	return nil
}

// CloseKafkaConsumer safely closes the Kafka consumer
func CloseKafkaConsumer(consumer sarama.Consumer) error {
	err := consumer.Close()
	if err != nil {
		log.Println("Error closing Kafka consumer:", err)
		return err
	}
	log.Println("Kafka consumer closed successfully")
	return nil
}

// PublishMessage publishes a message to a Kafka topic
func PublishMessage(producer sarama.SyncProducer, topic string, key string, value interface{}) error {
	// Serialize the message to JSON
	message, err := json.Marshal(value)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return err
	}

	// Create a Kafka message
	kafkaMessage := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(message),
	}

	// Send the message
	_, _, err = producer.SendMessage(kafkaMessage)
	if err != nil {
		log.Println("Error sending message:", err)
		return err
	}

	log.Printf("Message sent to topic %s: %s\n", topic, message)
	return nil
}

// ConsumeMessages starts consuming messages from the provided topic
func ConsumeMessages(consumer sarama.Consumer, topic string, eventHandler func(key, value []byte)) error {
	// Start a consumer for the topic
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		log.Println("Error fetching partitions:", err)
		return err
	}

	// Iterate over the partitions and consume messages
	for _, partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Println("Error starting partition consumer:", err)
			return err
		}

		// Consume messages from the partition
		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				// Handle the consumed message with the provided event handler
				eventHandler(msg.Key, msg.Value)
			}
		}(pc)
	}

	return nil
}
