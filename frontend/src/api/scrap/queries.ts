import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ICreateScrapPayload, IScrapsResponse } from "./types";
import { doGet, doPost } from "api/http";

const basePath = "/api/scrap";

export const useScraps = (userId: number | undefined) =>
  useQuery<IScrapsResponse>({
    queryKey: [`${basePath}/user/${userId}`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
    enabled: !!userId,
  });

export const useCreateScrap = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: ICreateScrapPayload) =>
      doPost({ url: `${basePath}`, body: payload }),
    onSettled: (_data, _error, variables) => {
      queryClient.invalidateQueries({
        queryKey: [`${basePath}/user/${variables.userRecipientId}`],
      });
    },
  });
};
