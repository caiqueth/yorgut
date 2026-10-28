import { useMutation, useQueryClient } from "@tanstack/react-query";
import { getFindMyWelcomeQueryKey } from "api/generated/endpoints/me/me";
import {
  createScrap,
  deleteScrap,
  getFindUserScrapsQueryKey,
  useFindUserScraps,
} from "api/generated/endpoints/scrap/scrap";
import { CreateScrapRequest } from "api/generated/model";
import { usePreviewParams } from "contexts/ProfilePreviewContext";

// One page of the user's scrapbook, newest first. Disabled until userId is
// known; the cast only satisfies the generated signature, the request never
// runs without an id.
// enabled false skips the request, e.g. for a scrapbook kept to friends.
export const useScraps = (userId: number | undefined, page: number, pageSize: number, enabled = true) =>
  useFindUserScraps(
    userId as number,
    { page, pageSize, ...usePreviewParams() },
    { query: { enabled: enabled && !!userId, placeholderData: (previous) => previous } },
  );

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

// recipientId is the owner of the scrapbook the scrap is in. Pending until
// the scrapbook reloads, so the deleted scrap doesn't linger on screen.
export const useDeleteScrap = (recipientId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (scrapId: number) => deleteScrap(scrapId),
    onSettled: () =>
      Promise.all([
        queryClient.invalidateQueries({ queryKey: getFindUserScrapsQueryKey(recipientId) }),
        queryClient.invalidateQueries({ queryKey: getFindMyWelcomeQueryKey() }),
      ]),
  });
};
