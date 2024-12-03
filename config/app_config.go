package config

import "os"

type AppConfig struct {
	port string
}

func loadAppConfig() AppConfig {
	return AppConfig{
		port: os.Getenv("APP_PORT"),
	}
}

func (ac AppConfig) Port() string {
	return ac.port
}
