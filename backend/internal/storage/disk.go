package storage

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Disk stores objects as files in one folder, which the API serves at the
// path of baseURL (see server.NewRouter).
type Disk struct {
	// Dir is absolute, so it doesn't depend on the working directory.
	Dir string
	// Path is baseURL's path, e.g. "/media".
	Path    string
	baseURL string
}

// NewDisk stores files in dir. baseURL must be absolute and have a path
// (http://localhost:8080/media), since the API serves the files there.
func NewDisk(dir string, baseURL string) (*Disk, error) {
	baseURL = strings.TrimRight(baseURL, "/")
	parsed, err := url.Parse(baseURL)
	if err != nil || !parsed.IsAbs() || parsed.Path == "" {
		return nil, fmt.Errorf("storage: base URL %q must be absolute and have a path, e.g. http://localhost:8080/media", baseURL)
	}

	absolute, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(absolute, 0o755); err != nil {
		return nil, err
	}
	return &Disk{Dir: absolute, Path: parsed.Path, baseURL: baseURL}, nil
}

func (d *Disk) Put(_ context.Context, key string, data []byte, _ string) error {
	if err := checkKey(key); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(d.Dir, key), data, 0o644)
}

func (d *Disk) Delete(_ context.Context, key string) error {
	if err := checkKey(key); err != nil {
		return err
	}
	if err := os.Remove(filepath.Join(d.Dir, key)); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return nil
}

func (d *Disk) URL(key string) string {
	return objectURL(d.baseURL, key)
}

func (d *Disk) String() string {
	return "folder " + d.Dir
}

// CanClear always allows it: the folder is on this machine.
func (d *Disk) CanClear() error {
	return nil
}

// Clear deletes every file in Dir, for the seed.
func (d *Disk) Clear(_ context.Context) error {
	if err := os.RemoveAll(d.Dir); err != nil {
		return err
	}
	return os.MkdirAll(d.Dir, 0o755)
}
