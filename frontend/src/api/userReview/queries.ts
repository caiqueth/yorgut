import { useMutation, useQueryClient } from "@tanstack/react-query";
import { getFindMyFriendsQueryKey, getFindMyWelcomeQueryKey } from "api/generated/endpoints/me/me";
import {
  getFindUserReviewsQueryKey,
  upsertUserReview,
  useFindUserReviews,
} from "api/generated/endpoints/user-review/user-review";
import { UpsertUserReviewRequest } from "api/generated/model";

export const useUserReviews = (userId: number) =>
  useFindUserReviews(userId, { query: { enabled: !!userId } });

export const useUpsertUserReview = (userId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: UpsertUserReviewRequest) => upsertUserReview(userId, payload),
    // The profile header's averages, and the friends list, which carries the
    // logged user's rating of each friend. Returning the promise keeps the
    // mutation pending until they reload, so the old rating doesn't flash back.
    onSettled: () =>
      Promise.all([
        queryClient.invalidateQueries({ queryKey: getFindUserReviewsQueryKey(userId) }),
        queryClient.invalidateQueries({ queryKey: getFindMyFriendsQueryKey() }),
        // The welcome box's "classifique seus amigos" step.
        queryClient.invalidateQueries({ queryKey: getFindMyWelcomeQueryKey() }),
      ]),
  });
};
