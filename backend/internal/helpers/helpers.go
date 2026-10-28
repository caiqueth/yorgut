package helpers

import (
	"io"
	"mime/multipart"
)

// FileToBytes reads a whole uploaded file.
func FileToBytes(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}
