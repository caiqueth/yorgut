package helpers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

var downloadClient = &http.Client{Timeout: 5 * time.Second}

// Download fetches url and returns its body, refusing bodies over maxBytes.
func Download(ctx context.Context, url string, maxBytes int64) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	response, err := downloadClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download %s: status %d", url, response.StatusCode)
	}
	// Read one byte past the limit to tell "exactly maxBytes" from "more".
	body, err := io.ReadAll(io.LimitReader(response.Body, maxBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(body)) > maxBytes {
		return nil, fmt.Errorf("download %s: larger than %d bytes", url, maxBytes)
	}
	return body, nil
}
