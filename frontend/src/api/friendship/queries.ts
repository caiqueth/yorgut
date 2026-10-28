import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  addFriend,
  removeFriend,
  replyFriendship,
} from "api/generated/endpoints/friendship/friendship";
import {
  getAmIFriendOfQueryKey,
  getFindMyBirthdaysQueryKey,
  getFindMyFriendsQueryKey,
  getFindMyPendingFriendshipsQueryKey,
  getFindMyWelcomeQueryKey,
} from "api/generated/endpoints/me/me";
import { CreateFriendshipRequest, ReplyFriendshipRequest } from "api/generated/model";

// Everything a friendship change (request, answer, removal) can change for the
// logged user: requests, friends, birthdays and the welcome stats.
const invalidateMyFriends = (queryClient: ReturnType<typeof useQueryClient>) => {
  for (const queryKey of [
    getFindMyPendingFriendshipsQueryKey(),
    getFindMyFriendsQueryKey(),
    getFindMyBirthdaysQueryKey(),
    getFindMyWelcomeQueryKey(),
  ]) {
    queryClient.invalidateQueries({ queryKey });
  }
};

export const useReplyFriendship = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: ReplyFriendshipRequest) => replyFriendship(payload),
    onSettled: () => {
      invalidateMyFriends(queryClient);
    },
  });
};

export const useAddFriend = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateFriendshipRequest) => addFriend(payload),
    onSettled: (_data, _error, variables) => {
      invalidateMyFriends(queryClient);
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
      invalidateMyFriends(queryClient);
      queryClient.invalidateQueries({ queryKey: getAmIFriendOfQueryKey(userId) });
    },
  });
};
