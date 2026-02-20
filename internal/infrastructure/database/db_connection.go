package database

import (
	"log"

	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	db, err := NewDBFromEnv()
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	return db
}
