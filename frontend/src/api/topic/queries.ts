import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { doGet, doPost } from "api/http";
import {
  ITopicCreateRequest,
  ITopicRepliesResponse,
  ITopicReplyRequest,
} from "./types";

const basePath = `${__API_BASE_URL__}/topics`;

export const useTopicReplies = (topicId: number) =>
  useQuery<ITopicRepliesResponse>({
    queryKey: [`${basePath}/${topicId}/replies`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useCreateTopic = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: ITopicCreateRequest) =>
      doPost({ url: `${basePath}`, body: payload }),
    onSettled: (_data, _error, variables) => {
      queryClient.invalidateQueries({
        queryKey: [`${__API_BASE_URL__}/communities/${variables.communityId}/hot-topics`],
      });
      queryClient.invalidateQueries({
        queryKey: [`${__API_BASE_URL__}/communities/${variables.communityId}/topics`],
      });
    },
  });
};

export const useReplyTopic = (topicId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: ITopicReplyRequest) =>
      doPost({ url: `${basePath}/${topicId}`, body: payload }),
    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: [`${basePath}/${topicId}/replies`],
      });
    },
  });
};
