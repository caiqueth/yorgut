package storage

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/backend/s3mem"
	"github.com/minio/minio-go/v7"
)

// fakeS3 starts an in-process S3 server with one bucket, "yorgut", and
// returns its endpoint (127.0.0.1:port).
func fakeS3(t *testing.T) string {
	t.Helper()
	backend := s3mem.New()
	if err := backend.CreateBucket("yorgut"); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(gofakes3.New(backend).Server())
	t.Cleanup(server.Close)
	return strings.TrimPrefix(server.URL, "http://")
}

func newTestS3(t *testing.T, endpoint string) *S3 {
	t.Helper()
	s, err := NewS3(context.Background(), S3Config{
		Endpoint:  endpoint,
		Region:    "us-east-1",
		Bucket:    "yorgut",
		AccessKey: "key",
		SecretKey: "secret",
		BaseURL:   "http://" + endpoint + "/yorgut/",
	})
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestS3(t *testing.T) {
	ctx := context.Background()
	endpoint := fakeS3(t)
	s := newTestS3(t, endpoint)

	t.Run("put, url and delete", func(t *testing.T) {
		key, err := PutImage(ctx, s, pngBytes(t))
		if err != nil {
			t.Fatal(err)
		}

		url := s.URL(key)
		if want := "http://" + endpoint + "/yorgut/" + key; url != want {
			t.Fatalf("URL = %q, want %q", url, want)
		}
		info, err := s.client.StatObject(ctx, "yorgut", key, minio.StatObjectOptions{})
		if err != nil {
			t.Fatal(err)
		}
		if info.ContentType != "image/png" {
			t.Fatalf("content type = %q, want image/png", info.ContentType)
		}
		// The fake serves path-style URLs like SeaweedFS does, so the URL loads.
		response, err := http.Get(url)
		if err != nil {
			t.Fatal(err)
		}
		body, _ := io.ReadAll(response.Body)
		response.Body.Close()
		if response.StatusCode != http.StatusOK || len(body) != len(pngBytes(t)) {
			t.Fatalf("GET %s: status %d, %d bytes", url, response.StatusCode, len(body))
		}

		if err := s.Delete(ctx, key); err != nil {
			t.Fatal(err)
		}
		if _, err := s.client.StatObject(ctx, "yorgut", key, minio.StatObjectOptions{}); err == nil {
			t.Fatal("object still exists after Delete")
		}
		if err := s.Delete(ctx, key); err != nil {
			t.Fatalf("deleting a missing key: %v", err)
		}
	})

	t.Run("invalid keys", func(t *testing.T) {
		for _, key := range []string{"", "..", "a/b.png", `a\b.png`} {
			if err := s.Put(ctx, key, []byte("x"), "image/png"); err == nil {
				t.Errorf("Put(%q) succeeded, want an error", key)
			}
		}
	})

	t.Run("clear", func(t *testing.T) {
		for range 3 {
			if _, err := PutImage(ctx, s, pngBytes(t)); err != nil {
				t.Fatal(err)
			}
		}
		if err := s.Clear(ctx); err != nil {
			t.Fatal(err)
		}
		for object := range s.client.ListObjects(ctx, "yorgut", minio.ListObjectsOptions{Recursive: true}) {
			t.Fatalf("object %q left after Clear (err %v)", object.Key, object.Err)
		}
	})

	t.Run("missing bucket", func(t *testing.T) {
		_, err := NewS3(ctx, S3Config{Endpoint: endpoint, Region: "us-east-1", Bucket: "other", BaseURL: "http://x/other"})
		if err == nil || !strings.Contains(err.Error(), "does not exist") {
			t.Fatalf("err = %v, want a missing bucket error", err)
		}
	})
}

func TestS3ClearOnlyOnLocalhost(t *testing.T) {
	for endpoint, local := range map[string]bool{
		"localhost:8333":                   true,
		"127.0.0.1:8333":                   true,
		"[::1]:8333":                       true,
		"localhost":                        true,
		"s3.us-east-1.amazonaws.com":       false,
		"abc.r2.cloudflarestorage.com:443": false,
		"localhost.example.com:8333":       false,
	} {
		s := &S3{endpoint: endpoint, bucket: "yorgut"}
		if got := s.CanClear() == nil; got != local {
			t.Errorf("CanClear on %s: allowed = %v, want %v", endpoint, got, local)
		}
	}
}

func TestNewFromEnv(t *testing.T) {
	t.Run("disk by default", func(t *testing.T) {
		t.Setenv("STORAGE", "")
		t.Setenv("UPLOADS_DIR", t.TempDir())
		t.Setenv("MEDIA_BASE_URL", "http://localhost:8080/media")
		store, err := NewFromEnv(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := store.(*Disk); !ok {
			t.Fatalf("store is %T, want *Disk", store)
		}
	})

	t.Run("s3", func(t *testing.T) {
		endpoint := fakeS3(t)
		t.Setenv("STORAGE", "s3")
		t.Setenv("S3_ENDPOINT", endpoint)
		t.Setenv("S3_USE_SSL", "false")
		t.Setenv("S3_REGION", "us-east-1")
		t.Setenv("S3_BUCKET", "yorgut")
		t.Setenv("S3_ACCESS_KEY", "key")
		t.Setenv("S3_SECRET_KEY", "secret")
		t.Setenv("MEDIA_BASE_URL", "http://"+endpoint+"/yorgut")
		store, err := NewFromEnv(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := store.(*S3); !ok {
			t.Fatalf("store is %T, want *S3", store)
		}
	})

	t.Run("errors", func(t *testing.T) {
		for name, env := range map[string]map[string]string{
			"no base URL":     {"STORAGE": "disk", "UPLOADS_DIR": t.TempDir(), "MEDIA_BASE_URL": ""},
			"no uploads dir":  {"STORAGE": "disk", "UPLOADS_DIR": "", "MEDIA_BASE_URL": "http://localhost:8080/media"},
			"unknown storage": {"STORAGE": "ftp", "MEDIA_BASE_URL": "http://localhost:8080/media"},
			"bad S3_USE_SSL":  {"STORAGE": "s3", "S3_USE_SSL": "maybe", "MEDIA_BASE_URL": "http://localhost/yorgut"},
		} {
			t.Run(name, func(t *testing.T) {
				for k, v := range env {
					t.Setenv(k, v)
				}
				if _, err := NewFromEnv(context.Background()); err == nil {
					t.Fatal("NewFromEnv succeeded, want an error")
				}
			})
		}
	})
}
