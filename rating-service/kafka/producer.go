package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/IBM/sarama"
)

var NotificationTopic = getNotificationTopic()
var producer sarama.SyncProducer

func getNotificationTopic() string {
	topic := os.Getenv("KAFKA_NOTIFICATION_TOPIC")
	if topic == "" {
		return "notification_events" // default
	}
	return topic
}

func InitKafkaProducer(brokers []string) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Net.SASL.Enable = true
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	config.Net.SASL.User = os.Getenv("KAFKA_USERNAME")
	config.Net.SASL.Password = os.Getenv("KAFKA_PASSWORD")
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = nil // uses default Root CAs

	config.Version = sarama.V2_5_0_0

	p, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return err
	}
	producer = p
	return nil
}

func SendNotification(ctx context.Context, payload interface{}) error {
	if producer == nil {
		log.Printf("[Kafka] Producer is not initialized")
		return nil
	}

	msgBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: NotificationTopic,
		Value: sarama.ByteEncoder(msgBytes),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		log.Printf("[Kafka] Failed to send message: %v", err)
		return err
	}

	log.Printf("[Kafka] Notification sent: %s", msgBytes)
	return nil
}
