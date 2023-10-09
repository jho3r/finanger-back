package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CorsMiddleware is the middleware to handle the cors.
func CorsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Check if the origin is allowed
		if !containsOrigin(origin, allowedOrigins) {
			c.AbortWithStatus(http.StatusForbidden)
		}

		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, X-API-KEY")
		c.Header("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)

			return
		}

		c.Next()
	}
}

// SplitOrigins splits the origins string into a slice.
func SplitOrigins(origins string) []string {
	return strings.Split(origins, ",")
}

// containsOrigin checks if the origin is allowed.
func containsOrigin(origin string, allowedOrigins []string) bool {
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}

	return false
}
