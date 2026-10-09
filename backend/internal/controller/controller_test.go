package controller

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"yorgut/internal/service"

	"github.com/gin-gonic/gin"
)

func TestRespondError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name   string
		err    error
		status int
	}{
		{"not found category", service.ErrNotFound, http.StatusNotFound},
		{"forbidden category", service.ErrForbidden, http.StatusForbidden},
		{"conflict category", service.ErrConflict, http.StatusConflict},
		{"invalid category", service.ErrInvalid, http.StatusBadRequest},
		{"specific not found", service.ErrTopicNotFound, http.StatusNotFound},
		{"specific conflict", service.ErrAlreadyMember, http.StatusConflict},
		{"specific invalid", service.ErrFirstNameRequired, http.StatusBadRequest},
		{"wrapped specific forbidden", fmt.Errorf("x: %w", service.ErrCommunityForbidden), http.StatusForbidden},
		{"invalid session", service.ErrInvalidSession, http.StatusUnauthorized},
		{"invalid google credential", service.ErrInvalidGoogleCredential, http.StatusUnauthorized},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)

			RespondError(ctx, tt.err)

			if rec.Code != tt.status {
				t.Fatalf("got status %d, want %d", rec.Code, tt.status)
			}
		})
	}

	t.Run("unknown error is a generic 500 and does not leak", func(t *testing.T) {
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)

		RespondError(ctx, errors.New("pq: password authentication failed for user secret"))

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("got status %d, want 500", rec.Code)
		}
		body := rec.Body.String()
		if strings.Contains(body, "password") || strings.Contains(body, "secret") {
			t.Fatalf("raw error leaked in body: %s", body)
		}
		if got, want := strings.TrimSpace(body), `{"error":"internal server error"}`; got != want {
			t.Fatalf("got body %s, want %s", got, want)
		}
	})
}
