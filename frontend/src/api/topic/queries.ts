import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  getFindCommunityHotTopicsQueryKey,
  getFindCommunityTopicsQueryKey,
} from "api/generated/endpoints/community/community";
import {
  createTopic,
  getFindTopicRepliesQueryKey,
  replyTopic,
  useFindTopicReplies,
} from "api/generated/endpoints/topics/topics";
import { CreateTopicRequest, ReplyTopicRequest } from "api/generated/model";

export const useTopicReplies = (topicId: number) => useFindTopicReplies(topicId);

export const useCreateTopic = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateTopicRequest) => createTopic(payload),
    onSettled: (_data, _error, variables) => {
      queryClient.invalidateQueries({
        queryKey: getFindCommunityHotTopicsQueryKey(variables.communityId),
      });
      queryClient.invalidateQueries({
        queryKey: getFindCommunityTopicsQueryKey(variables.communityId),
      });
    },
  });
};

export const useReplyTopic = (topicId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: ReplyTopicRequest) => replyTopic(topicId, payload),
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: getFindTopicRepliesQueryKey(topicId) });
    },
  });
};
