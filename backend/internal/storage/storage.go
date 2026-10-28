// Package storage keeps uploaded files (pictures) out of the API process.
// The database stores only an object's key ("<uuid>.jpg"); the ObjectStore
// knows where the bytes live and which URL serves them.
package storage

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type ObjectStore interface {
	// Put stores data under key, replacing any object already there.
	Put(ctx context.Context, key string, data []byte, contentType string) error
	// Delete removes the object; deleting a missing key is not an error.
	Delete(ctx context.Context, key string) error
	// URL is the absolute URL clients load the object from.
	URL(key string) string
}

var extensions = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
}

// PutImage stores a PNG or JPEG under a new random key and returns the key.
// Validating the picture (size, dimensions) is the caller's job.
func PutImage(ctx context.Context, store ObjectStore, data []byte) (string, error) {
	contentType := http.DetectContentType(data)
	extension, ok := extensions[contentType]
	if !ok {
		return "", fmt.Errorf("storage: unsupported image type %q", contentType)
	}

	key := uuid.NewString() + extension
	if err := store.Put(ctx, key, data, contentType); err != nil {
		return "", err
	}
	return key, nil
}

// checkKey keeps every key a plain name: no folders, so a key can't point
// outside the Disk folder, and both stores accept the same keys.
func checkKey(key string) error {
	if key == "" || key == "." || key == ".." || key != filepath.Base(key) || strings.ContainsAny(key, `/\`) {
		return fmt.Errorf("storage: invalid key %q", key)
	}
	return nil
}

func objectURL(baseURL string, key string) string {
	return baseURL + "/" + url.PathEscape(key)
}

// NewFromEnv builds the store the API and the seed use. STORAGE picks it:
//
//   - "disk" (default): UPLOADS_DIR is the folder, MEDIA_BASE_URL the URL the
//     API serves it at (http://localhost:8080/media).
//   - "s3": S3_ENDPOINT, S3_USE_SSL, S3_REGION, S3_BUCKET, S3_ACCESS_KEY,
//     S3_SECRET_KEY, and MEDIA_BASE_URL is the bucket's public URL.
func NewFromEnv(ctx context.Context) (ObjectStore, error) {
	baseURL := os.Getenv("MEDIA_BASE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("MEDIA_BASE_URL must be set")
	}

	switch kind := os.Getenv("STORAGE"); kind {
	case "", "disk":
		dir := os.Getenv("UPLOADS_DIR")
		if dir == "" {
			return nil, fmt.Errorf("UPLOADS_DIR must be set when STORAGE is disk")
		}
		return NewDisk(dir, baseURL)
	case "s3":
		return NewS3FromEnv(ctx)
	default:
		return nil, fmt.Errorf("STORAGE must be disk or s3, not %q", kind)
	}
}

// NewS3FromEnv builds the S3 store from the S3_* variables and
// MEDIA_BASE_URL, whatever STORAGE says.
func NewS3FromEnv(ctx context.Context) (*S3, error) {
	useSSL := true
	if v := os.Getenv("S3_USE_SSL"); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("S3_USE_SSL: %w", err)
		}
		useSSL = parsed
	}
	return NewS3(ctx, S3Config{
		Endpoint:  os.Getenv("S3_ENDPOINT"),
		UseSSL:    useSSL,
		Region:    os.Getenv("S3_REGION"),
		Bucket:    os.Getenv("S3_BUCKET"),
		AccessKey: os.Getenv("S3_ACCESS_KEY"),
		SecretKey: os.Getenv("S3_SECRET_KEY"),
		BaseURL:   os.Getenv("MEDIA_BASE_URL"),
	})
}

// CopyDir puts every file directly in dir into store, keyed by its file name
// (the keys the database already holds), and returns how many it copied.
// Putting an existing key again overwrites it, so it can be rerun.
func CopyDir(ctx context.Context, dir string, store ObjectStore) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	copied := 0
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return copied, err
		}
		if err := store.Put(ctx, entry.Name(), data, http.DetectContentType(data)); err != nil {
			return copied, fmt.Errorf("copying %s: %w", entry.Name(), err)
		}
		copied++
	}
	return copied, nil
}
