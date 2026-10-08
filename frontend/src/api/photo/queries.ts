import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { doDelete, doGet, doPost } from "api/http";
import { IPhotosResponse, IUploadPhotoPayload } from "./types";

const basePath = `${__API_BASE_URL__}/photos`;

export const usePhotos = (userId: number) =>
  useQuery<IPhotosResponse>({
    queryKey: [`${basePath}/user/${userId}`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
    enabled: !!userId,
  });

// userId is the logged user's id (the album owner), to refresh their album.
export const useUploadPhoto = (userId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ picture, caption }: IUploadPhotoPayload) => {
      const form = new FormData();
      form.append("picture", picture);
      form.append("caption", caption);
      return doPost({ url: `${basePath}`, body: form });
    },
    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: [`${basePath}/user/${userId}`],
      });
    },
  });
};

// userId is the logged user's id (the album owner), to refresh their album.
export const useDeletePhoto = (userId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (photoId: number) => doDelete({ url: `${basePath}/${photoId}` }),
    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: [`${basePath}/user/${userId}`],
      });
    },
  });
};
