package infrastructure

import (
	"geodishdiscoveryservice/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Infrastructure struct {
	awsConfig   aws.Config
	db          *mongo.Client
	kafka       *KafkaResources
}

func NewInfrastructure(cfg *config.Config) (*Infrastructure, error) {
	awsConfig, err := initAWS(cfg.AWSConfig())
	if err != nil {
		return nil, err
	}

	db, err := initDB(cfg.DatabaseConfig())
	if err != nil {
		return nil, err
	}

	kafkaResources, err := initKafka(cfg.KafkaConfig())
	if err != nil {
		return nil, err
	}

	return &Infrastructure{
		awsConfig:   awsConfig,
		db:          db,
		kafka:     kafkaResources,
	}, nil
}

// Getter for awsConfig
func (i *Infrastructure) AWSConfig() aws.Config {
	return i.awsConfig
}

// Getter for db
func (i *Infrastructure) DB() *mongo.Client {
	return i.db
}

func (i *Infrastructure) KafkaResources() *KafkaResources {
	return i.kafka
}
