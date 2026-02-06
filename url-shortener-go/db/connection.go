package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"url-shortener-go/models"
)

var DB *gorm.DB

// Connect connects to the PostgreSQL database and returns a *gorm.DB object.
// It also automatically migrates the models to the database.
// The database connection string is hardcoded to use the 'DB_HOST', 'DB_USER', 'DB_PASSWORD', 'DB_NAME', and 'DB_PORT' environment variables.
// The SSL mode is disabled.
// If the connection fails, it returns a non-nil error.
// If the connection is successful and the migration is also successful, it returns nil.
func Connect() error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Auto-migrate models
	return DB.AutoMigrate(&models.URL{})
}
