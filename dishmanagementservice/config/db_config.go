package config

import "os"

type DatabaseConfig struct {
	driver   string
	host     string
	port     string
	dbname   string
	username string
	password string
}

func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		driver:   os.Getenv("DB_DRIVER"),
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
		dbname:   os.Getenv("DB_NAME"),
		username: os.Getenv("DB_USER"),
		password: os.Getenv("DB_PASSWORD"),
	}
}

func (dc DatabaseConfig) Driver() string {
	return dc.driver
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
