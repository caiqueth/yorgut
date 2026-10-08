import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { doDelete, doGet, doPut } from "api/http";
import { IUserFanSummaryResponse } from "./types";

const basePath = `${__API_BASE_URL__}/users`;

export const useUserFans = (userId: number) =>
  useQuery<IUserFanSummaryResponse>({
    queryKey: [`${basePath}/${userId}/fans`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
    enabled: !!userId,
  });

// Pass true to become a fan of the user, false to remove the star.
export const useSetFan = (userId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (isFan: boolean) =>
      isFan
        ? doPut({ url: `${basePath}/${userId}/fan`, body: {} })
        : doDelete({ url: `${basePath}/${userId}/fan` }),
    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: [`${basePath}/${userId}/fans`],
      });
    },
  });
};
