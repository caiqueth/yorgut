import { useMutation, useQueryClient } from "@tanstack/react-query";
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
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: getFindUserReviewsQueryKey(userId) });
    },
  });
};
