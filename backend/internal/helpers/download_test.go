package helpers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDownload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/small":
			w.Write([]byte("12345"))
		case "/big":
			w.Write([]byte(strings.Repeat("x", 11)))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	ctx := context.Background()

	body, err := Download(ctx, server.URL+"/small", 10)
	if err != nil || string(body) != "12345" {
		t.Fatalf("small: got %q, %v", body, err)
	}
	if _, err := Download(ctx, server.URL+"/big", 10); err == nil {
		t.Error("big: got no error, want one over the limit")
	}
	if _, err := Download(ctx, server.URL+"/missing", 10); err == nil {
		t.Error("missing: got no error, want one for the 404")
	}
}
