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
		if got, want := strings.TrimSpace(body), `{"error":"Algo deu errado do nosso lado. Tente de novo mais tarde."}`; got != want {
			t.Fatalf("got body %s, want %s", got, want)
		}
	})
}

func TestRespondErrorShowsOnlyTheUserMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			"wrapping context stays out",
			fmt.Errorf("loading album 3: %w", service.ErrPhotoNotFound),
			service.ErrPhotoNotFound.Error(),
		},
		{
			"Google's reason stays out",
			fmt.Errorf("%w: idtoken: token expired", service.ErrInvalidGoogleCredential),
			service.ErrInvalidGoogleCredential.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)

			RespondError(ctx, tt.err)

			if got, want := strings.TrimSpace(rec.Body.String()), fmt.Sprintf(`{"error":%q}`, tt.want); got != want {
				t.Fatalf("got body %s, want %s", got, want)
			}
		})
	}
}

func TestIdParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		value  string
		wantId uint
		wantOk bool
	}{
		{"42", 42, true},
		{"0", 0, false},
		{"-1", 0, false},
		{"0x1f", 0, false},
		{"abc", 0, false},
		{"", 0, false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q", tt.value), func(t *testing.T) {
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Params = gin.Params{{Key: "userId", Value: tt.value}}

			id, ok := idParam(ctx, "userId")
			if id != tt.wantId || ok != tt.wantOk {
				t.Fatalf("got (%d, %v), want (%d, %v)", id, ok, tt.wantId, tt.wantOk)
			}
			if !ok && rec.Code != http.StatusBadRequest {
				t.Fatalf("got status %d, want 400", rec.Code)
			}
		})
	}
}
