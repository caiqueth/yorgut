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

// The topic and one page of its replies, oldest first. enabled false skips
// the request, e.g. in a private community's forum when the logged user isn't
// a member (the API would answer 403).
export const useTopicReplies = (topicId: number, page: number, pageSize: number, enabled = true) =>
  useFindTopicReplies(
    topicId,
    { page, pageSize },
    { query: { enabled, placeholderData: (previous) => previous } },
  );

export const useCreateTopic = () => {
  const queryClient = useQueryClient();
  return useMutation({
    meta: { showsErrorInline: true },
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

// communityId refreshes the forum lists, which show reply counts.
export const useReplyTopic = (topicId: number, communityId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: ReplyTopicRequest) => replyTopic(topicId, payload),
    onSettled: () => {
      for (const queryKey of [
        getFindTopicRepliesQueryKey(topicId),
        getFindCommunityTopicsQueryKey(communityId),
        getFindCommunityHotTopicsQueryKey(communityId),
      ]) {
        queryClient.invalidateQueries({ queryKey });
      }
    },
  });
};
