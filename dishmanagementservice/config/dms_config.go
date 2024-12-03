package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	appConfig      AppConfig
	databaseConfig DatabaseConfig
	kafkaConfig    KafkaConfig
	awsConfig      AWSConfig
}

var dishManagementServiceConfig *Config

func init() {
	dishManagementServiceConfig = loadConfig()
}

func loadConfig() *Config {
	// Load .env if it exists
	if err := godotenv.Load("dms.env"); err != nil {
		log.Println("No dms.env file found:", err)
	}

	return &Config{
		appConfig:      loadAppConfig(),
		databaseConfig: loadDatabaseConfig(),
		kafkaConfig:    loadKafkaConfig(),
		awsConfig:      loadAWSConfig(),
	}
}

func GetDishManagementServiceConfig() *Config {
	return dishManagementServiceConfig
}

func (cfg *Config) AppConfig() AppConfig {
	return cfg.appConfig
}

func (cfg *Config) DatabaseConfig() DatabaseConfig {
	return cfg.databaseConfig
}

func (cfg *Config) KafkaConfig() KafkaConfig {
	return cfg.kafkaConfig
}

func (cfg *Config) AWSConfig() AWSConfig {
	return cfg.awsConfig
}
