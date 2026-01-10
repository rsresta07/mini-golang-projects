package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers the routes for the web application.
// It takes a pointer to a *gin.Engine and a pointer to a *gorm.DB as arguments.
// It registers a temporary health check at the "/health" endpoint, which returns a JSON response with the message "ok".
// It also calls registerAuthRoutes to register the authentication routes.
func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	// temporary health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	registerAuthRoutes(r, db)
}
