package server

import (
	"math/rand"
	"net/http"
	"strings"
	"time"
	"yorgut/internal/controller"
	"yorgut/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks the Yorgut session token from the Authorization
// header and stores the logged user in the context as "user".
func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
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

		user, err := authService.UserFromSession(parts[1])
		if err != nil {
			controller.RespondError(c, err)
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func DelayMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		delay := time.Duration(rand.Intn(10)+1) * time.Second
		time.Sleep(delay)
		c.Next()
	}
}
