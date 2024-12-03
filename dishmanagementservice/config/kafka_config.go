package config

import "os"

type KafkaConfig struct {
	host    			string
	port    			string
	topicDishCreated	string
	topicDishUpdated	string
	topicDishDeleted	string
	topicOrderCreated	string
	topicOrderUpdated	string
	topicOrderDeleted	string
	groupID				string
}

func loadKafkaConfig() KafkaConfig {
	return KafkaConfig{
		host:    			os.Getenv("KAFKA_HOST"),
		port:    			os.Getenv("KAFKA_PORT"),
		topicDishCreated:   os.Getenv("KAFKA_TOPIC_DISH_CREATED"),
		topicDishUpdated:   os.Getenv("KAFKA_TOPIC_DISH_UPDATED"),
		topicDishDeleted:   os.Getenv("KAFKA_TOPIC_DISH_DELETED"),
		topicOrderCreated:   os.Getenv("KAFKA_TOPIC_ORDER_CREATED"),
		topicOrderUpdated:   os.Getenv("KAFKA_TOPIC_ORDER_UPDATED"),
		topicOrderDeleted:   os.Getenv("KAFKA_TOPIC_ORDER_DELETED"),
		groupID: 			os.Getenv("KAFKA_GROUP_ID"),
	}
}

func (kc KafkaConfig) Host() string {
	return kc.host
}

func (kc KafkaConfig) Port() string {
	return kc.port
}

func (kc KafkaConfig) TopicDishCreated() string {
	return kc.topicDishCreated
}

func (kc KafkaConfig) TopicDishUpdated() string {
	return kc.topicDishUpdated
}

func (kc KafkaConfig) TopicDishDeleted() string {
	return kc.topicDishDeleted
}

func (kc KafkaConfig) TopicOrderCreated() string {
	return kc.topicOrderCreated
}

func (kc KafkaConfig) TopicOrderUpdated() string {
	return kc.topicOrderUpdated
}

func (kc KafkaConfig) TopicOrderDeleted() string {
	return kc.topicOrderDeleted
}

func (kc KafkaConfig) GroupID() string {
	return kc.groupID
}
