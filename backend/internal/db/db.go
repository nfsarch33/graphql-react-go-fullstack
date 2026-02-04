package db

import (
	"log"

	"github.com/nfsarch33/graphql-react-go-fullstack/backend/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDB initializes SQLite database and runs migrations
func InitDB() (*gorm.DB, error) {
	// Connect to SQLite database (creates if doesn't exist)
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")

	// Run migrations
	err = db.AutoMigrate(&models.Todo{})
	if err != nil {
		return nil, err
	}

	log.Println("Database migrations completed")

	return db, nil
}
