import { useMutation, useQueryClient } from "@tanstack/react-query";
import { doDelete, doPost } from "api/http";
import { IAddFriendRequest, IFriendshipReplyRequest } from "./types";

const basePathFriendship = `${__API_BASE_URL__}/friendships`;
const basePathMe = `${__API_BASE_URL__}/me`;

export const useReplyFriendship = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: IFriendshipReplyRequest) =>
      doPost({ url: `${basePathFriendship}/reply`, body: payload }),
    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: [`${basePathMe}/friends/pending`],
      });
      queryClient.invalidateQueries({
        queryKey: [`${basePathMe}/friends`],
      });
    },
  });
};

export const useAddFriend = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: IAddFriendRequest) =>
      doPost({ url: `${basePathFriendship}`, body: payload }),
    onSettled: (_data, _error, variables) => {
      queryClient.invalidateQueries({
        queryKey: [`${basePathMe}/friends/pending`],
      });
      queryClient.invalidateQueries({
        queryKey: [`${basePathMe}/friends`],
      });
      queryClient.invalidateQueries({
        queryKey: [`${basePathMe}/friends/${variables.userRequesteeId}`],
      });
    },
  });
};

export const useRemoveFriend = (userId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => doDelete({ url: `${basePathFriendship}/friend/${userId}` }),
    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: [`${basePathMe}/friends/pending`],
      });
      queryClient.invalidateQueries({
        queryKey: [`${basePathMe}/friends`],
      });
      queryClient.invalidateQueries({
        queryKey: [`${basePathMe}/friends/${userId}`],
      });
    },
  });
};

export const useDeleteFriendship = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (friendshipId: number) =>
      doDelete({ url: `${basePathFriendship}/${friendshipId}` }),
    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: [`${basePathMe}/friends/pending`],
      });
      queryClient.invalidateQueries({
        queryKey: [`${basePathMe}/friends`],
      });
    },
  });
};
