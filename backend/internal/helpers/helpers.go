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

func AddRootToPath(relativePath string) string {
	pwd, _ := os.Getwd()
	return path.Join(pwd, relativePath)
}

func PathOrAbsolute(dir string, absolute bool) string {
	if absolute {
		return AddRootToPath(dir)
	}
	return dir
}

func GetStoragePath(absolute bool) string {
	return PathOrAbsolute("storage", absolute)
}

func GetUploadsPath(absolute bool) string {
	return PathOrAbsolute(path.Join("storage", "uploads"), absolute)
}

func MustParseUint(s string) uint {
	val, _ := strconv.ParseUint(s, 10, 64)
	return uint(val)
}

func MustParseInt(s string) int {
	val, _ := strconv.ParseInt(s, 10, 64)
	return int(val)
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

func UploadFileToLocalStorage(file []byte) (string, error) {
	uploadsPathAbs := GetUploadsPath(true)
	uploadsPathRel := GetUploadsPath(false)

	if err := os.MkdirAll(uploadsPathAbs, os.ModePerm); err != nil {
		return "", err
	}

	fileName := uuid.New().String()

	fileLocation := path.Join(uploadsPathAbs, fileName)
	fileNameLocationRel := path.Join(uploadsPathRel, fileName)

	err := os.WriteFile(fileLocation, file, 0644)
	if err != nil {
		return "", err
	}

	return fileNameLocationRel, nil
}

func UploadPicture(picture []byte) (*string, error) {
	imageHasValidExtension := CheckValidImage(picture)
	if !imageHasValidExtension {
		return nil, errors.New("invalid image format")
	}

	storedFileName, err := UploadFileToLocalStorage(picture)
	if err != nil {
		return nil, err
	}
	return &storedFileName, nil
}
