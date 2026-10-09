import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  becomeFan,
  getFindUserFansQueryKey,
  stopBeingFan,
  useFindUserFans,
} from "api/generated/endpoints/user-fan/user-fan";

export const useUserFans = (userId: number) =>
  useFindUserFans(userId, { query: { enabled: !!userId } });

// Pass true to become a fan of the user, false to remove the star.
export const useSetFan = (userId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (isFan: boolean) => (isFan ? becomeFan(userId) : stopBeingFan(userId)),
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: getFindUserFansQueryKey(userId) });
    },
  });
};
