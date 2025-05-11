import { useMutation, useQueryClient } from "@tanstack/react-query";
import { doDelete, doPost } from "api/http";
import { IAddFriendRequest, IFriendshipReplyRequest } from "./types";

export const basePath = "/api/friendship";

export const useReplyFriendship = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: IFriendshipReplyRequest) =>
      doPost({ url: `${basePath}/reply`, body: payload }),
    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: [`/api/me/friends/pending`],
      });
      queryClient.invalidateQueries({
        queryKey: [`/api/me/friends`],
      });
    },
  });
};

export const useAddFriend = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: IAddFriendRequest) =>
      doPost({ url: `${basePath}`, body: payload }),
    onSettled: (_data, _error, variables) => {
      queryClient.invalidateQueries({
        queryKey: [`/api/me/friends/pending`],
      });
      queryClient.invalidateQueries({
        queryKey: [`/api/me/friends`],
      });
      queryClient.invalidateQueries({
        queryKey: [`/api/me/friends/${variables.userRequesteeId}`],
      });
    },
  });
};

export const useRemoveFriend = (userId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => doDelete({ url: `${basePath}/friend/${userId}` }),
    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: [`/api/me/friends/pending`],
      });
      queryClient.invalidateQueries({
        queryKey: [`/api/me/friends`],
      });
      queryClient.invalidateQueries({
        queryKey: [`/api/me/friends/${userId}`],
      });
    },
  });
};

export const useDeleteFriendship = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (friendshipId: number) =>
      doDelete({ url: `${basePath}/${friendshipId}` }),
    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: [`/api/me/friends/pending`],
      });
      queryClient.invalidateQueries({
        queryKey: [`/api/me/friends`],
      });
    },
  });
};
