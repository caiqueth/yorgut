package dto

type CreatePhotoRequest struct {
	PhotoAlbumID uint
	URL          string
	Caption      string
}

type FindPhotoOfUserResponse struct{}

type FindPhotoOfAlbumResponse struct{}
