import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  closePoll,
  createPoll,
  getFindCommunityPollsQueryKey,
  getFindPollQueryKey,
  useFindCommunityPolls,
  useFindPoll,
  votePoll,
} from "api/generated/endpoints/poll/poll";
import { CreatePollRequest } from "api/generated/model";

// One page of the community's polls, newest first. enabled false skips the
// request, e.g. in a private community when the logged user isn't a member
// (the API would answer 403).
export const useCommunityPolls = (communityId: number, page: number, pageSize: number, enabled = true) =>
  useFindCommunityPolls(
    communityId,
    { page, pageSize },
    { query: { enabled, placeholderData: (previous) => previous } },
  );

// The poll with its options; each option's votes are null until the logged
// user voted or the poll closed.
export const usePoll = (pollId: number, enabled = true) =>
  useFindPoll(pollId, { query: { enabled } });

export const useCreatePoll = (communityId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    meta: { showsErrorInline: true },
    mutationFn: (payload: CreatePollRequest) => createPoll(communityId, payload),
    onSettled: () =>
      queryClient.invalidateQueries({ queryKey: getFindCommunityPollsQueryKey(communityId) }),
  });
};

// Voting and closing change the poll and its line in the community's list
// (vote count, closed).
const invalidatePoll = (queryClient: ReturnType<typeof useQueryClient>, pollId: number, communityId: number) => {
  queryClient.invalidateQueries({ queryKey: getFindPollQueryKey(pollId) });
  queryClient.invalidateQueries({ queryKey: getFindCommunityPollsQueryKey(communityId) });
};

// mutate(optionId) votes for that option.
export const useVotePoll = (pollId: number, communityId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (optionId: number) => votePoll(pollId, { optionId }),
    onSettled: () => invalidatePoll(queryClient, pollId, communityId),
  });
};

export const useClosePoll = (pollId: number, communityId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => closePoll(pollId),
    onSettled: () => invalidatePoll(queryClient, pollId, communityId),
  });
};
