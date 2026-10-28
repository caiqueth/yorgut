// uploads-to-s3 copies the pictures of a disk store (UPLOADS_DIR) into the S3
// bucket from the S3_* variables, keeping their keys, so switching STORAGE
// from disk to s3 keeps every picture. Safe to rerun.
//
//	go run ./cmd/uploads-to-s3
package main

import (
	"context"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"yorgut/internal/storage"
)

func main() {
	ctx := context.Background()

	dir := os.Getenv("UPLOADS_DIR")
	if dir == "" {
		log.Fatal("❌ UPLOADS_DIR must be set: the folder to copy from")
	}
	bucket, err := storage.NewS3FromEnv(ctx)
	if err != nil {
		log.Fatalf("❌ %v", err)
	}

	copied, err := storage.CopyDir(ctx, dir, bucket)
	if err != nil {
		log.Fatalf("❌ copied %d files from %s to %v, then: %v", copied, dir, bucket, err)
	}
	log.Printf("✅ copied %d files from %s to %v", copied, dir, bucket)
}
