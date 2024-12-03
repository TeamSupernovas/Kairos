package infrastructure

import (
	"database/sql"
	"dishmanagementservice/config"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type Infrastructure struct {
	awsConfig   aws.Config
	db          *sql.DB
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
func (i *Infrastructure) DB() *sql.DB {
	return i.db
}

// Getter for kafkaWriter
func (i *Infrastructure) KafkaResources() *KafkaResources {
	return i.kafka
}
