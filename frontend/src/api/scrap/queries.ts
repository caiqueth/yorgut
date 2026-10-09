import { useMutation, useQueryClient } from "@tanstack/react-query";
import { getFindMyWelcomeQueryKey } from "api/generated/endpoints/me/me";
import {
  createScrap,
  getFindUserScrapsQueryKey,
  useFindUserScraps,
} from "api/generated/endpoints/scrap/scrap";
import { CreateScrapRequest } from "api/generated/model";

// Disabled until userId is known; the cast only satisfies the generated
// signature, the request never runs without an id.
export const useScraps = (userId: number | undefined) =>
  useFindUserScraps(userId as number, { query: { enabled: !!userId } });

export const useCreateScrap = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateScrapRequest) => createScrap(payload),
    onSettled: (_data, _error, variables) => {
      queryClient.invalidateQueries({
        queryKey: getFindUserScrapsQueryKey(variables.userRecipientId),
      });
      queryClient.invalidateQueries({ queryKey: getFindMyWelcomeQueryKey() });
    },
  });
};
