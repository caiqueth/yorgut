package helpers

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/google/uuid"
)

func GetStoragePath() string {
	pwd, _ := os.Getwd()
	return path.Join(pwd, "storage")
}

func GetUploadsPath() string {
	storagePath := GetStoragePath()
	return path.Join(storagePath, "uploads")
}

func MustParseUint(s string) uint {
	val, _ := strconv.ParseUint(s, 10, 64)
	return uint(val)
}

func FileToBytes(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func CheckValidImage(d []byte) bool {
	validExtensions := map[string]bool{
		"image/png":  true,
		"image/jpg":  true,
		"image/jpeg": true,
	}
	contentType := http.DetectContentType(d)

	return validExtensions[contentType]
}

func WriteFileToLocalStorage(subDir string, file []byte) (string, error) {
	storagePath := GetStoragePath()
	fileName := path.Join(subDir, uuid.New().String())

	err := os.WriteFile(path.Join(storagePath, fileName), file, 0644)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
