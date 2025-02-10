package database

import (
	"log"

	"github.com/oragazzo/todo_grpc_tls/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect establishes a connection to the database and performs migrations
func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Println("Database connection successful")

	// Auto migrate the models
	if err := db.AutoMigrate(&models.Todo{}); err != nil {
		return nil, err
	}

	return db, nil
}
