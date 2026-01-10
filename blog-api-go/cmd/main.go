package main

import (
	"blog-api-go/config"
	"blog-api-go/routes"
	"github.com/gin-gonic/gin"
)

// main is the entry point for the application.
// It connects to the database, initializes the routes
// and starts the web server.
func main() {
	db := config.ConnectDB()

	r := gin.Default()
	routes.RegisterRoutes(r, db)
	r.Run(":8080")
}
