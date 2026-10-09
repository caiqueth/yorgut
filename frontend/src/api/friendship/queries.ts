import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  addFriend,
  deleteFriendship,
  removeFriend,
  replyFriendship,
} from "api/generated/endpoints/friendship/friendship";
import {
  getAmIFriendOfQueryKey,
  getFindMyFriendsQueryKey,
  getFindMyPendingFriendshipsQueryKey,
} from "api/generated/endpoints/me/me";
import { CreateFriendshipRequest, ReplyFriendshipRequest } from "api/generated/model";

export const useReplyFriendship = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: ReplyFriendshipRequest) => replyFriendship(payload),
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: getFindMyPendingFriendshipsQueryKey() });
      queryClient.invalidateQueries({ queryKey: getFindMyFriendsQueryKey() });
    },
  });
};

export const useAddFriend = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateFriendshipRequest) => addFriend(payload),
    onSettled: (_data, _error, variables) => {
      queryClient.invalidateQueries({ queryKey: getFindMyPendingFriendshipsQueryKey() });
      queryClient.invalidateQueries({ queryKey: getFindMyFriendsQueryKey() });
      queryClient.invalidateQueries({
        queryKey: getAmIFriendOfQueryKey(variables.userRequesteeId),
      });
    },
  });
};

export const useRemoveFriend = (userId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => removeFriend(userId),
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: getFindMyPendingFriendshipsQueryKey() });
      queryClient.invalidateQueries({ queryKey: getFindMyFriendsQueryKey() });
      queryClient.invalidateQueries({ queryKey: getAmIFriendOfQueryKey(userId) });
    },
  });
};

export const useDeleteFriendship = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (friendshipId: number) => deleteFriendship(friendshipId),
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: getFindMyPendingFriendshipsQueryKey() });
      queryClient.invalidateQueries({ queryKey: getFindMyFriendsQueryKey() });
    },
  });
};
