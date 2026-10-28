package server

import (
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
	"yorgut/internal/controller"
	"yorgut/internal/service"
	"yorgut/internal/transport"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// AuthMiddleware checks the Yorgut session token from the Authorization
// header and stores the logged user in the context as "user".
func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, transport.ErrorResponse{Error: service.ErrInvalidSession.Error()})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, transport.ErrorResponse{Error: service.ErrInvalidSession.Error()})
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

// delayFromEnv reads API_DELAY, a development aid that slows every API
// response to show loading states: a duration ("800ms", "1s") or a random
// range ("300ms-2s"). ok is false when it's unset; an invalid value is logged
// and ignored.
func delayFromEnv() (shortest time.Duration, longest time.Duration, ok bool) {
	value := strings.TrimSpace(os.Getenv("API_DELAY"))
	if value == "" {
		return 0, 0, false
	}

	low, high, isRange := strings.Cut(value, "-")
	shortest, errMin := time.ParseDuration(low)
	longest, errMax := shortest, error(nil)
	if isRange {
		longest, errMax = time.ParseDuration(high)
	}
	if errMin != nil || errMax != nil || shortest < 0 || longest < shortest {
		log.Warn().Str("API_DELAY", value).Msg(`ignoring API_DELAY: use a duration like "800ms" or a range like "300ms-2s"`)
		return 0, 0, false
	}
	return shortest, longest, true
}

// DelayMiddleware waits a random time between shortest and longest before handling
// each request. See delayFromEnv.
func DelayMiddleware(shortest time.Duration, longest time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		delay := shortest
		if longest > shortest {
			delay += time.Duration(rand.Int63n(int64(longest - shortest)))
		}
		time.Sleep(delay)
		c.Next()
	}
}
