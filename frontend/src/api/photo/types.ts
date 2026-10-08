export interface IPhoto {
  id: number;
  url: string;
  caption: string;
  createdAt: Date;
}

export interface IPhotosResponse {
  photos: IPhoto[];
  count: number;
  // How many photos the personal album can hold.
  maxPhotos: number;
  // True when the logged user owns the album.
  canUpload: boolean;
}

export interface IUploadPhotoPayload {
  picture: File;
  caption: string;
}
