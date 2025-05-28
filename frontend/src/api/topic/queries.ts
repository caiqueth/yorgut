import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { doGet, doPost } from "api/http";
import {
  ITopicCreateRequest,
  ITopicRepliesResponse,
  ITopicReplyRequest,
} from "./types";

const basePath = "/api/topic";

export const useTopicReplies = (topicId: number) =>
  useQuery<ITopicRepliesResponse>({
    queryKey: [`${basePath}/${topicId}/replies`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
    enabled: !!topicId,
  });

export const useCreateTopic = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: ITopicCreateRequest) =>
      doPost({ url: `${basePath}`, body: payload }),
    onSettled: (_data, _error, variables) => {
      queryClient.invalidateQueries({
        queryKey: [`/api/community/${variables.communityId}/hot_topics`],
      });
      queryClient.invalidateQueries({
        queryKey: [`/api/community/${variables.communityId}/topics`],
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
        queryKey: [`/api/topic/${topicId}/replies`],
      });
    },
  });
};
