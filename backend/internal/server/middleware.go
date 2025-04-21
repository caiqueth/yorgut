package server

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var ContextTokenKey = os.Getenv("CONTEXT_TOKEN_KEY")

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		token := parts[1]

		// Do something to validate the token...

		// Store the token in the context
		c.Set(ContextTokenKey, token)

		c.Next()
	}
}
