package service

import (
	"context"
	"yorgut/internal/helpers"
	"yorgut/internal/storage"

	"github.com/rs/zerolog/log"
)

// Limits for every uploaded picture: album photos, profile and community
// pictures. Uploads may be big; what we store is shrunk (see preparePicture).
const (
	MaxUploadBytes  = 5 << 20    // 5 MB, what a user may send
	MaxUploadPixels = 50_000_000 // 50 megapixels, so decoding can't exhaust memory
	MaxPictureBytes = 1 << 20    // 1 MB, what we store
	MaxPictureSide  = 1024       // pixels, the longest side of what we store
)

var (
	ErrInvalidPicture  = newError(ErrInvalid, "A foto precisa ser PNG ou JPEG.")
	ErrPictureTooLarge = newError(ErrInvalid, "A foto pode ter no máximo 5 MB.")
	ErrPictureTooBig   = newError(ErrInvalid, "A foto pode ter no máximo 50 megapixels.")
)

// preparePicture checks an upload (PNG or JPEG, at most 5 MB and 50
// megapixels) and returns the picture to store: at most 1024 pixels on its
// longest side and 1 MB, rotated upright and without its metadata (see
// helpers.ShrinkImage).
func preparePicture(picture []byte) ([]byte, error) {
	if len(picture) > MaxUploadBytes {
		return nil, ErrPictureTooLarge
	}

	// Only reads the header: PNG and JPEG are the formats it knows.
	width, height, err := helpers.ImageDimensions(picture)
	if err != nil {
		return nil, ErrInvalidPicture
	}
	if width*height > MaxUploadPixels {
		return nil, ErrPictureTooBig
	}

	shrunk, err := helpers.ShrinkImage(picture, MaxPictureSide, MaxPictureBytes)
	if err != nil {
		return nil, ErrInvalidPicture
	}
	return shrunk, nil
}

// deletePicture removes the file of a picture nothing points to anymore.
// Google profile pictures are stored as full URLs, not keys, so they are
// skipped. The database change already happened, so a failure is only
// logged: it leaves an orphan file, not a broken picture.
func deletePicture(ctx context.Context, store storage.ObjectStore, picture *string) {
	if picture == nil || *picture == "" || isExternalPicture(*picture) {
		return
	}
	if err := store.Delete(ctx, *picture); err != nil {
		log.Warn().Err(err).Str("key", *picture).Msg("could not delete a picture that is no longer used")
	}
}
