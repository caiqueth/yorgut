import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  deletePhoto,
  getFindUserPhotosQueryKey,
  uploadPhoto,
  useFindUserPhotos,
} from "api/generated/endpoints/photo/photo";
import { UploadPhotoBody } from "api/generated/model";
import { usePreviewParams } from "contexts/ProfilePreviewContext";

// enabled false skips the request, e.g. for an album kept to friends.
export const usePhotos = (userId: number, enabled = true) =>
  useFindUserPhotos(userId, usePreviewParams(), { query: { enabled: enabled && !!userId } });

// userId is the logged user's id (the album owner), to refresh their album.
export const useUploadPhoto = (userId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    meta: { showsErrorInline: true },
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
