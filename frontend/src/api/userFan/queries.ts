import { useMutation, useQueryClient } from "@tanstack/react-query";
import { getFindMyFriendsQueryKey } from "api/generated/endpoints/me/me";
import {
  becomeFan,
  getFindUserFanListQueryKey,
  getFindUserFansQueryKey,
  stopBeingFan,
  useFindUserFanList,
  useFindUserFans,
} from "api/generated/endpoints/user-fan/user-fan";

export const useUserFans = (userId: number) =>
  useFindUserFans(userId, { query: { enabled: !!userId } });

// One page of the user's fans, the newest first.
export const useUserFanList = (userId: number, page: number, pageSize: number) =>
  useFindUserFanList(
    userId,
    { page, pageSize },
    { query: { enabled: !!userId, placeholderData: (previous) => previous } },
  );

// Pass true to become a fan of the user, false to remove the star.
export const useSetFan = (userId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (isFan: boolean) => (isFan ? becomeFan(userId) : stopBeingFan(userId)),
    // The friends list carries the logged user's fan star on each friend.
    // Returning the promise keeps the mutation pending until everything
    // reloads, so the old star doesn't flash back.
    onSettled: () =>
      Promise.all([
        queryClient.invalidateQueries({ queryKey: getFindUserFansQueryKey(userId) }),
        queryClient.invalidateQueries({ queryKey: getFindUserFanListQueryKey(userId) }),
        queryClient.invalidateQueries({ queryKey: getFindMyFriendsQueryKey() }),
      ]),
  });
};
