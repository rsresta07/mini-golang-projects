package routes

import (
	"blog-api-go/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// registerAuthRoutes registers the authentication routes for the web application.
// It takes a pointer to a *gin.Engine and a pointer to a *gorm.DB as arguments.
// It registers a new instance of AuthController with the provided DB and sets up the routes for registering a new user and logging in an existing user.
func registerAuthRoutes(r *gin.Engine, db *gorm.DB) {
	auth := controllers.NewAuthController(db)
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", auth.Register)
		authRoutes.POST("/login", auth.Login)
	}
}
