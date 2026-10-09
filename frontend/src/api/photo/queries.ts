import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  deletePhoto,
  getFindUserPhotosQueryKey,
  uploadPhoto,
  useFindUserPhotos,
} from "api/generated/endpoints/photo/photo";
import { UploadPhotoBody } from "api/generated/model";

export const usePhotos = (userId: number) =>
  useFindUserPhotos(userId, { query: { enabled: !!userId } });

// userId is the logged user's id (the album owner), to refresh their album.
export const useUploadPhoto = (userId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: UploadPhotoBody) => uploadPhoto(payload),
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: getFindUserPhotosQueryKey(userId) });
    },
  });
};

// userId is the logged user's id (the album owner), to refresh their album.
export const useDeletePhoto = (userId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (photoId: number) => deletePhoto(photoId),
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: getFindUserPhotosQueryKey(userId) });
    },
  });
};
