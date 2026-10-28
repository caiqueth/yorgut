package storage

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func pngBytes(t *testing.T) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestDisk(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	disk, err := NewDisk(dir, "http://localhost:8080/media/")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("put, url and delete", func(t *testing.T) {
		if err := disk.Put(ctx, "a.png", []byte("x"), "image/png"); err != nil {
			t.Fatal(err)
		}
		if got, _ := os.ReadFile(filepath.Join(dir, "a.png")); string(got) != "x" {
			t.Fatalf("file content = %q, want %q", got, "x")
		}
		if got, want := disk.URL("a.png"), "http://localhost:8080/media/a.png"; got != want {
			t.Fatalf("URL = %q, want %q", got, want)
		}
		if disk.Path != "/media" {
			t.Fatalf("Path = %q, want /media", disk.Path)
		}

		if err := disk.Delete(ctx, "a.png"); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(filepath.Join(dir, "a.png")); !os.IsNotExist(err) {
			t.Fatalf("file still exists after Delete (err %v)", err)
		}
		if err := disk.Delete(ctx, "a.png"); err != nil {
			t.Fatalf("deleting a missing key: %v", err)
		}
	})

	t.Run("keys can't leave the folder", func(t *testing.T) {
		for _, key := range []string{"", ".", "..", "../a.png", "sub/a.png", `sub\a.png`} {
			if filepath.Separator == '/' && key == `sub\a.png` {
				continue // a valid file name outside Windows
			}
			if err := disk.Put(ctx, key, []byte("x"), "image/png"); err == nil {
				t.Errorf("Put(%q) succeeded, want an error", key)
			}
		}
	})

	t.Run("base URL needs a path", func(t *testing.T) {
		for _, base := range []string{"http://localhost:8080", "http://localhost:8080/", "/media", ""} {
			if _, err := NewDisk(dir, base); err == nil {
				t.Errorf("NewDisk(%q) succeeded, want an error", base)
			}
		}
	})
}

func TestPutImage(t *testing.T) {
	ctx := context.Background()
	store := NewMemory()

	key, err := PutImage(ctx, store, pngBytes(t))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(key, ".png") || !store.Has(key) {
		t.Fatalf("key %q: want a stored .png key", key)
	}

	if _, err := PutImage(ctx, store, []byte("not an image")); err == nil {
		t.Fatal("PutImage accepted a non-image")
	}
	if store.Len() != 1 {
		t.Fatalf("store has %d objects, want 1", store.Len())
	}
}

func TestCopyDir(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	picture := pngBytes(t)
	for _, name := range []string{"old-key-without-extension", "new.png"} {
		if err := os.WriteFile(filepath.Join(dir, name), picture, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}

	store := NewMemory()
	for run := 1; run <= 2; run++ { // rerunning copies the same keys again
		copied, err := CopyDir(ctx, dir, store)
		if err != nil {
			t.Fatal(err)
		}
		if copied != 2 || store.Len() != 2 || !store.Has("old-key-without-extension") || !store.Has("new.png") {
			t.Fatalf("run %d: copied %d, store has %d objects, want the 2 files by name", run, copied, store.Len())
		}
	}
}
