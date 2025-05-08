package config

import (
	"fmt"
	"os"
)

type DatabaseConfig struct {
	connectionUri string
	host     string
	port     string
	dbname   string
	username string
	password string
}

func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		connectionUri: os.Getenv("MONGO_URI"),
		host:     os.Getenv("MONGO_HOST"),
		port:     os.Getenv("MONGO_PORT"),
		dbname:   os.Getenv("MONGO_DB_NAME"),
		username: os.Getenv("MONGO_USERNAME"),
		password: os.Getenv("MONGO_PASSWORD"),
	}
}

func (dc DatabaseConfig) ConnectionURI() string {
	if dc.connectionUri != "" {
		return dc.connectionUri
	}
	var uri string
    if dc.Username() != "" && dc.Password() != "" {
        // Use authentication if credentials are provided
        uri = fmt.Sprintf("mongodb://%s:%s@%s:%s", dc.Username(), dc.Password(), dc.Host(), dc.Port())
    } else {
        // No authentication
        uri = fmt.Sprintf("mongodb://%s:%s", dc.Host(), dc.Port())
    }

	return uri
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
