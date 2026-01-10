package config

import (
	"blog-api-go/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB connects to the PostgreSQL database and returns a *gorm.DB object.
// It also automatically migrates the User and Blog models to the database.
// The database connection string is hardcoded to use the 'blog' database on localhost with the
// username 'postgres' and password 'password'. The SSL mode is disabled.
// You may want to change this to use environment variables or a configuration file in a production setting.
func ConnectDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=password dbname=blog port=5432 sslmode=disable"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&models.User{}, &models.Blog{})
	return db
}
