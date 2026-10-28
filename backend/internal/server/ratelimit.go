package server

import (
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"
	"yorgut/internal/model"
	"yorgut/internal/transport"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// The request limits. Each is a token bucket: Burst requests at once, then one
// more every Every. Reads aren't limited.
var (
	// Logins, per client address: each one checks a Google token.
	loginLimit = limit{Every: 6 * time.Second, Burst: 10}
	// Picture uploads, per user: each is resized and compressed.
	uploadLimit = limit{Every: 12 * time.Second, Burst: 5}
	// Every other change (scraps, ratings, friend requests...), per user.
	// Generous: rating friends on /friends is one request per click.
	writeLimit = limit{Every: time.Second, Burst: 60}
)

type limit struct {
	Every time.Duration
	Burst int
}

// keyedLimiter keeps one token bucket per key (a user id, a client address)
// and forgets the keys that have been idle for a while.
type keyedLimiter struct {
	limit     limit
	mu        sync.Mutex
	buckets   map[string]*bucket
	lastSweep time.Time
}

type bucket struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// idleAfter is how long an unused bucket is kept. By then it has refilled,
// so forgetting it changes nothing.
const idleAfter = 10 * time.Minute

func newKeyedLimiter(l limit) *keyedLimiter {
	return &keyedLimiter{limit: l, buckets: map[string]*bucket{}}
}

// allow takes a token from key's bucket. When there is none, it returns how
// long until there will be.
func (k *keyedLimiter) allow(key string, now time.Time) (bool, time.Duration) {
	k.mu.Lock()
	defer k.mu.Unlock()

	if now.Sub(k.lastSweep) > idleAfter {
		for key, b := range k.buckets {
			if now.Sub(b.lastSeen) > idleAfter {
				delete(k.buckets, key)
			}
		}
		k.lastSweep = now
	}

	b, ok := k.buckets[key]
	if !ok {
		b = &bucket{limiter: rate.NewLimiter(rate.Every(k.limit.Every), k.limit.Burst)}
		k.buckets[key] = b
	}
	b.lastSeen = now

	reservation := b.limiter.ReserveN(now, 1)
	if delay := reservation.DelayFrom(now); delay > 0 {
		reservation.CancelAt(now)
		return false, delay
	}
	return true, 0
}

// RateLimit answers 429, with a Retry-After header, to requests over the
// limit. key picks whose bucket a request uses; an empty key isn't limited.
func RateLimit(l limit, key func(*gin.Context) string) gin.HandlerFunc {
	limiter := newKeyedLimiter(l)
	return func(c *gin.Context) {
		k := key(c)
		if k == "" {
			c.Next()
			return
		}
		if ok, wait := limiter.allow(k, time.Now()); !ok {
			c.Header("Retry-After", strconv.Itoa(int(math.Ceil(wait.Seconds()))))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, transport.ErrorResponse{
				Error: "Muitas ações em pouco tempo. Espere um pouco e tente de novo.",
			})
			return
		}
		c.Next()
	}
}

// byClientAddress keys requests by the visitor's address. Behind Cloudflare
// that is the CF-Connecting-IP header (see TrustedPlatform in NewRouter).
func byClientAddress(c *gin.Context) string {
	return c.ClientIP()
}

// byUser keys requests by the logged user (set by AuthMiddleware).
func byUser(c *gin.Context) string {
	user, ok := c.Get("user")
	if !ok {
		return ""
	}
	return strconv.FormatUint(uint64(user.(*model.User).ID), 10)
}

// writesByUser is byUser for requests that change something; reads pass.
func writesByUser(c *gin.Context) string {
	switch c.Request.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return ""
	}
	return byUser(c)
}
