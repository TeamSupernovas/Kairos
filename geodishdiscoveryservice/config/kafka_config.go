package config

import "os"

type KafkaConfig struct {
	host    			string
	port    			string
	username 			string
	password 			string
	topicDishCreated	string
	topicDishUpdated	string
	topicDishDeleted	string
	groupID				string
}

func loadKafkaConfig() KafkaConfig {
	return KafkaConfig{
		host:    			os.Getenv("KAFKA_HOST"),
		port:    			os.Getenv("KAFKA_PORT"),
		username: 			os.Getenv("KAFKA_USERNAME"),
		password: 			os.Getenv("KAFKA_PASSWORD"),
		topicDishCreated:   os.Getenv("KAFKA_TOPIC_DISH_CREATED"),
		topicDishUpdated:   os.Getenv("KAFKA_TOPIC_DISH_UPDATED"),
		topicDishDeleted:   os.Getenv("KAFKA_TOPIC_DISH_DELETED"),
		groupID: 			os.Getenv("KAFKA_GROUP_ID"),
	}
}

func (kc KafkaConfig) Host() string {
	return kc.host
}

func (kc KafkaConfig) Port() string {
	return kc.port
}

func (kc KafkaConfig) Username() string {
	return kc.username
}

func (kc KafkaConfig) Password() string {
	return kc.password
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

func (kc KafkaConfig) GroupID() string {
	return kc.groupID
}
