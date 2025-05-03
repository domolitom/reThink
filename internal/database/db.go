package database

import (
	"fmt"
	"log"
	"os"

	"github.com/domolitom/reThink/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the database connection
var DB *gorm.DB

// Connect establishes a connection to the database and performs migrations
func Connect() {
	var err error

	// Get database credentials from environment variables
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "prediction_social"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	// Create the connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port)

	// Connect to the database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to database successfully")

	// Auto-migrate the schema
	err = DB.AutoMigrate(&models.User{}, &models.Market{}, &models.Prediction{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migration completed")
}
