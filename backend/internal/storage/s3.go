package storage

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3 stores objects in a bucket of any S3-compatible service (SeaweedFS,
// Cloudflare R2, Backblaze B2, AWS). Clients load them from baseURL, which
// points at the bucket (or at a proxy or CDN in front of it); the bucket must
// allow public reads.
type S3 struct {
	client   *minio.Client
	bucket   string
	endpoint string
	baseURL  string
}

type S3Config struct {
	Endpoint  string // host[:port], without scheme, e.g. "localhost:8333"
	UseSSL    bool
	Region    string // optional for most services
	Bucket    string
	AccessKey string
	SecretKey string
	BaseURL   string // public URL of the bucket, e.g. "http://localhost:8333/yorgut"
}

// NewS3 connects to the service and checks that the bucket exists; creating
// it is a deploy step, not the API's job.
func NewS3(ctx context.Context, cfg S3Config) (*S3, error) {
	if cfg.Endpoint == "" || cfg.Bucket == "" || cfg.BaseURL == "" {
		return nil, fmt.Errorf("storage: S3 needs an endpoint, a bucket and a base URL")
	}
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, err
	}

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("storage: checking bucket %q at %s: %w", cfg.Bucket, cfg.Endpoint, err)
	}
	if !exists {
		return nil, fmt.Errorf("storage: bucket %q does not exist at %s", cfg.Bucket, cfg.Endpoint)
	}

	return &S3{
		client:   client,
		bucket:   cfg.Bucket,
		endpoint: cfg.Endpoint,
		baseURL:  strings.TrimRight(cfg.BaseURL, "/"),
	}, nil
}

func (s *S3) Put(ctx context.Context, key string, data []byte, contentType string) error {
	if err := checkKey(key); err != nil {
		return err
	}
	_, err := s.client.PutObject(ctx, s.bucket, key, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
		// Keys are random and never reused, so the content never changes.
		CacheControl: "public, max-age=31536000, immutable",
	})
	return err
}

func (s *S3) Delete(ctx context.Context, key string) error {
	if err := checkKey(key); err != nil {
		return err
	}
	return s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{})
}

func (s *S3) URL(key string) string {
	return objectURL(s.baseURL, key)
}

func (s *S3) String() string {
	return fmt.Sprintf("bucket %q at %s", s.bucket, s.endpoint)
}

// CanClear refuses endpoints that aren't on this machine, so the seed can't
// empty a real bucket.
func (s *S3) CanClear() error {
	if !isLocalEndpoint(s.endpoint) {
		return fmt.Errorf("storage: refusing to clear bucket %q on %s, which is not localhost", s.bucket, s.endpoint)
	}
	return nil
}

// Clear deletes every object in the bucket, for the seed (see CanClear).
func (s *S3) Clear(ctx context.Context) error {
	if err := s.CanClear(); err != nil {
		return err
	}
	// RemoveObjects skips listing errors, so pass the objects on and keep the
	// first error.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var listErr error
	listed := make(chan struct{})
	objects := make(chan minio.ObjectInfo)
	go func() {
		defer close(listed)
		defer close(objects)
		for object := range s.client.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{Recursive: true}) {
			if object.Err != nil {
				if listErr == nil {
					listErr = object.Err
				}
				continue
			}
			select {
			case objects <- object:
			case <-ctx.Done():
				return
			}
		}
	}()

	var removeErr error
	for result := range s.client.RemoveObjects(ctx, s.bucket, objects, minio.RemoveObjectsOptions{}) {
		if result.Err != nil && removeErr == nil {
			removeErr = result.Err
		}
	}
	cancel()
	<-listed
	if listErr != nil {
		return listErr
	}
	return removeErr
}

func isLocalEndpoint(endpoint string) bool {
	host, _, err := net.SplitHostPort(endpoint)
	if err != nil {
		host = endpoint
	}
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}
