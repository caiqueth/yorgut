package dto

import "yorgut/internal/model"

type PhotoAlbumDto struct {
	UserID uint
	Name   string
}

type CreatePhotoAlbumRequest struct {
	UserID uint
	Name   string
}

type FindPhotoAlbumOfUserResponse struct {
	PhotoAlbums []PhotoAlbumDto `json:"photoAlbums"`
}

func (r *FindPhotoAlbumOfUserResponse) FromPhotoAlbums(photoAlbums []model.PhotoAlbum) {
	r.PhotoAlbums = make([]PhotoAlbumDto, len(photoAlbums))

	for i, pa := range photoAlbums {
		r.PhotoAlbums[i] = PhotoAlbumDto{
			Name:   pa.Name,
			UserID: pa.UserID,
		}
	}
}
