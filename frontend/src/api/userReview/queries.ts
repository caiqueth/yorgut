import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { doGet, doPut } from "api/http";
import { IUpsertUserReviewPayload, IUserReviewSummaryResponse } from "./types";

const basePath = `${__API_BASE_URL__}/users`;

export const useUserReviews = (userId: number) =>
  useQuery<IUserReviewSummaryResponse>({
    queryKey: [`${basePath}/${userId}/reviews`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
    enabled: !!userId,
  });

export const useUpsertUserReview = (userId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: IUpsertUserReviewPayload) =>
      doPut({ url: `${basePath}/${userId}/review`, body: payload }),
    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: [`${basePath}/${userId}/reviews`],
      });
    },
  });
};
