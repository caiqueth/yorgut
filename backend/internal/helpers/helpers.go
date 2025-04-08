package helpers

import (
	"errors"
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

func ReadLocalFile(path string) ([]byte, error) {
	dat, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return dat, nil
}

func UploadFileToLocalStorage(subDir string, file []byte) (string, error) {
	uploadsPath := GetUploadsPath()
	fileDir := path.Join(uploadsPath, subDir)

	if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
		return "", err
	}

	fileName := path.Join(fileDir, uuid.New().String())

	err := os.WriteFile(fileName, file, 0644)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func UploadPicture(picture []byte, subPath string) (*string, error) {
	imageHasValidExtension := CheckValidImage(picture)
	if !imageHasValidExtension {
		return nil, errors.New("invalid image format")
	}

	storedFileName, err := UploadFileToLocalStorage(subPath, picture)
	if err != nil {
		return nil, err
	}
	return &storedFileName, nil
}

func UploadProfilePicture(picture []byte) (*string, error) {
	return UploadPicture(picture, "profile_pictures")
}

func UploadCommunityPicture(picture []byte) (*string, error) {
	return UploadPicture(picture, "community_pictures")
}
