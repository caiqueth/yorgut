package server

import (
	"net/http"
	"strings"
	"yorgut/internal/helpers"
	"yorgut/internal/service"

	"github.com/gin-gonic/gin"
)

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

		// Parse the Token
		parsed, err := helpers.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Could not parse the token"})
			return
		}

		// Store the token in the context
		helpers.SetJWTToken(c, parsed)

		c.Next()
	}
}

func InjectUserMiddleware(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := helpers.GetUserEmail(c)

		user, err := userService.FindUserByEmail(email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User fetch failed"})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
