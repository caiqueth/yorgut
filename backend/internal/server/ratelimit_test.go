package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"yorgut/internal/model"

	"github.com/gin-gonic/gin"
)

func TestKeyedLimiter(t *testing.T) {
	limiter := newKeyedLimiter(limit{Every: 10 * time.Second, Burst: 3})
	now := time.Now()

	for i := 1; i <= 3; i++ {
		if ok, _ := limiter.allow("ana", now); !ok {
			t.Fatalf("request %d within the burst was refused", i)
		}
	}
	ok, wait := limiter.allow("ana", now)
	if ok || wait <= 0 || wait > 10*time.Second {
		t.Fatalf("over the burst: got ok %v, wait %v; want refused with a wait up to 10s", ok, wait)
	}
	if ok, _ := limiter.allow("bia", now); !ok {
		t.Error("another key shares ana's bucket")
	}
	if ok, _ := limiter.allow("ana", now.Add(10*time.Second)); !ok {
		t.Error("a token wasn't back after Every")
	}

	// Idle buckets are forgotten.
	limiter.allow("caio", now.Add(time.Minute))
	limiter.allow("caio", now.Add(time.Minute+idleAfter+time.Second))
	if _, kept := limiter.buckets["ana"]; kept {
		t.Error("an idle bucket was kept")
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	user := &model.User{}
	user.ID = 1
	router.Use(func(c *gin.Context) { c.Set("user", user) })
	router.Use(RateLimit(limit{Every: time.Minute, Burst: 1}, writesByUser))
	router.GET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.POST("/x", func(c *gin.Context) { c.Status(http.StatusOK) })

	do := func(method string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(method, "/x", nil))
		return rec
	}

	if rec := do(http.MethodPost); rec.Code != http.StatusOK {
		t.Fatalf("first write: got %d, want 200", rec.Code)
	}
	rec := do(http.MethodPost)
	if rec.Code != http.StatusTooManyRequests || rec.Header().Get("Retry-After") == "" {
		t.Fatalf("second write: got %d, Retry-After %q; want 429 with Retry-After", rec.Code, rec.Header().Get("Retry-After"))
	}
	for i := 0; i < 5; i++ {
		if rec := do(http.MethodGet); rec.Code != http.StatusOK {
			t.Fatalf("read %d: got %d, want 200 (reads aren't limited)", i, rec.Code)
		}
	}
}
