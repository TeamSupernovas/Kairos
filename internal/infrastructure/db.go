package infrastructure

import (
	"database/sql"
	"dishmanagementservice/config"
	"fmt"
	"log"
)

func initDB(dbConfig config.DatabaseConfig) (*sql.DB, error) {

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host(),
		dbConfig.Port(),
		dbConfig.Username(),
		dbConfig.Password(),
		dbConfig.DBName(),
	)
	
	db, err := sql.Open(dbConfig.Driver(), connectionString)

	if (err != nil) {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	// Test the connection 
	err = db.Ping()
	if (err != nil) {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	log.Println("Database connected successfully")
	return db,nil
}