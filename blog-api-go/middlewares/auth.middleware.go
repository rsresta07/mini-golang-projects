package middlewares

import (
	"net/http"
	"strings"

	"blog-api-go/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware that checks if a valid JWT token is provided in the Authorization header.
// If the token is missing, invalid, or expired, it returns an error response with a 401 status.
// If the token is valid, it sets the user ID in the request context and calls the next handler in the chain.
// It can be used to protect routes that require authentication.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			c.Abort()
			return
		}

		userID, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
