package config

import "os"

type DatabaseConfig struct {
	host     string
	port     string
	dbname   string
	username string
	password string
}

func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		host:     os.Getenv("MONGO_HOST"),
		port:     os.Getenv("MONGO_PORT"),
		dbname:   os.Getenv("MONGO_DB_NAME"),
		username: os.Getenv("MONGO_USERNAME"),
		password: os.Getenv("MONGO_PASSWORD"),
	}
}

func (dc DatabaseConfig) Host() string {
	return dc.host
}

func (dc DatabaseConfig) Port() string {
	return dc.port
}

func (dc DatabaseConfig) DBName() string {
	return dc.dbname
}

func (dc DatabaseConfig) Username() string {
	return dc.username
}

func (dc DatabaseConfig) Password() string {
	return dc.password
}
