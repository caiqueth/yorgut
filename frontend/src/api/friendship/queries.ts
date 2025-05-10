import { useMutation, useQueryClient } from "@tanstack/react-query";
import { doPost } from "api/http";
import { IFriendshipReplyRequest } from "./types";

export const basePath = "/api/friendship";

export const useReplyFriendship = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: IFriendshipReplyRequest) =>
      doPost({ url: `${basePath}/reply`, body: payload }),
    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: [`/api/me/friends/pending`],
      });
      queryClient.invalidateQueries({
        queryKey: [`/api/me/friends`],
      });
    },
  });
};
