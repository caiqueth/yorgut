package service

import (
	"fmt"
	"yorgut/internal/helpers"
)

// Limits for every uploaded picture: album photos and profile pictures.
const (
	MaxPictureBytes     = 1 << 20 // 1 MB
	MaxPictureDimension = 1024    // pixels, for both width and height
)

var (
	ErrInvalidPicture  = newError(ErrInvalid, "the picture must be a PNG or JPEG image")
	ErrPictureTooLarge = newError(ErrInvalid, "the picture must be at most 1 MB")
	ErrPictureTooBig   = newError(ErrInvalid, fmt.Sprintf("the picture must be at most %dx%d pixels", MaxPictureDimension, MaxPictureDimension))
)

// validatePicture enforces the upload limits: PNG or JPEG, at most 1 MB and
// 1024x1024 pixels.
func validatePicture(picture []byte) error {
	if len(picture) > MaxPictureBytes {
		return ErrPictureTooLarge
	}
	if !helpers.CheckValidImage(picture) {
		return ErrInvalidPicture
	}

	width, height, err := helpers.ImageDimensions(picture)
	if err != nil {
		return ErrInvalidPicture
	}
	if width > MaxPictureDimension || height > MaxPictureDimension {
		return ErrPictureTooBig
	}

	return nil
}
