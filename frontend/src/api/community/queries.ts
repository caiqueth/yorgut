import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  createCommunity,
  getFindCommunityJoinRequestsQueryKey,
  getFindCommunityMembersQueryKey,
  getFindCommunityProfileQueryKey,
  getFindCommunityQueryKey,
  getFindInvitableFriendsQueryKey,
  inviteToCommunity,
  getSearchCommunitiesQueryKey,
  replyCommunityJoinRequest,
  updateCommunity,
  updateCommunityPicture,
  useFindCommunity,
  useFindCommunityHotTopics,
  useFindCommunityJoinRequests,
  useFindCommunityMembers,
  useFindCommunityProfile,
  useFindCommunityTopics,
  useFindInvitableFriends,
  useFindRelatedCommunities,
  useSearchCommunities as useSearchCommunitiesQuery,
} from "api/generated/endpoints/community/community";
import {
  getFindMyCommunitiesQueryKey,
  getFindMyMembershipsQueryKey,
  getFindMyWelcomeQueryKey,
} from "api/generated/endpoints/me/me";
import { CreateCommunityBody, UpdateCommunityRequest } from "api/generated/model";

export const useCommunity = (communityId: number) => useFindCommunity(communityId);

// One page of the community's members, most recently seen first.
export const useCommunityMembers = (communityId: number, page: number, pageSize: number) =>
  useFindCommunityMembers(communityId, { page, pageSize }, { query: { placeholderData: (previous) => previous } });

export const useCommunityProfile = (communityId: number) =>
  useFindCommunityProfile(communityId);

// 9 random other communities (the API picks the same ones all day, so reloads
// don't reshuffle).
export const useCommunityRelated = (communityId: number) =>
  useFindRelatedCommunities(communityId);

// Everything that shows a community's name, category, description, picture
// or privacy. Making a private community public also accepts its join
// requests, which changes its members.
const communityKeys = (communityId: number) => [
  getFindCommunityQueryKey(communityId),
  getFindCommunityJoinRequestsQueryKey(communityId),
  getFindCommunityMembersQueryKey(communityId),
  getFindCommunityProfileQueryKey(communityId),
  getSearchCommunitiesQueryKey(),
  getFindMyCommunitiesQueryKey(),
  getFindMyMembershipsQueryKey(),
];

// Owner only; the backend answers 403 to anyone else.
export const useUpdateCommunity = (communityId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    meta: { showsErrorInline: true },
    mutationFn: (payload: UpdateCommunityRequest) => updateCommunity(communityId, payload),
    onSettled: () => {
      for (const queryKey of communityKeys(communityId)) {
        queryClient.invalidateQueries({ queryKey });
      }
    },
  });
};

// Owner only; the backend answers 403 to anyone else.
export const useUpdateCommunityPicture = (communityId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    meta: { showsErrorInline: true },
    mutationFn: (picture: File) => updateCommunityPicture(communityId, { picture }),
    onSettled: () => {
      for (const queryKey of communityKeys(communityId)) {
        queryClient.invalidateQueries({ queryKey });
      }
    },
  });
};

// Sent as multipart/form-data so the optional picture travels with the fields.
export const useCreateCommunity = () => {
  const queryClient = useQueryClient();
  return useMutation({
    meta: { showsErrorInline: true },
    mutationFn: (payload: CreateCommunityBody) => createCommunity(payload),
    onSettled: () => {
      for (const queryKey of [
        getFindMyCommunitiesQueryKey(),
        getFindMyMembershipsQueryKey(),
        getSearchCommunitiesQueryKey(),
        // The welcome box counts communities for its next steps.
        getFindMyWelcomeQueryKey(),
      ]) {
        queryClient.invalidateQueries({ queryKey });
      }
    },
  });
};

// enabled false skips the request, e.g. for a private community's forum when
// the logged user isn't a member (the API would answer 403).
export const useCommunityHotTopics = (communityId: number, enabled = true) =>
  useFindCommunityHotTopics(communityId, { query: { enabled } });

// One page of the community's topics, newest first.
export const useCommunityTopics = (communityId: number, page: number, pageSize: number, enabled = true) =>
  useFindCommunityTopics(
    communityId,
    { page, pageSize },
    { query: { enabled, placeholderData: (previous) => previous } },
  );

// Who asked to join a private community. Owner only (403 for anyone else),
// so pass enabled false for other users.
export const useCommunityJoinRequests = (communityId: number, enabled: boolean) =>
  useFindCommunityJoinRequests(communityId, { query: { enabled } });

// Owner only. Accepting makes the user a member.
export const useReplyJoinRequest = (communityId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ userId, accept }: { userId: number; accept: boolean }) =>
      replyCommunityJoinRequest(communityId, userId, { accept }),
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: getFindCommunityJoinRequestsQueryKey(communityId) });
      queryClient.invalidateQueries({ queryKey: getFindCommunityMembersQueryKey(communityId) });
    },
  });
};

// Communities the logged user hasn't joined, filtered by name (empty = all).
// The key starts with the search path, so joining/leaving can refresh every
// search at once by invalidating getSearchCommunitiesQueryKey().
export const useSearchCommunities = (name: string, page: number, pageSize: number) =>
  useSearchCommunitiesQuery(
    { name, page, pageSize },
    { query: { placeholderData: (previous) => previous } },
  );

// One page of the logged user's friends, each saying whether they are already
// a member of the community or already invited. Members only (403 for anyone
// else), so pass enabled false for other users.
export const useInvitableFriends = (communityId: number, page: number, pageSize: number, enabled: boolean) =>
  useFindInvitableFriends(
    communityId,
    { page, pageSize },
    { query: { enabled, placeholderData: (previous) => previous } },
  );

// mutate(userId) invites that friend to the community.
export const useInviteToCommunity = (communityId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (userId: number) => inviteToCommunity(communityId, { userId }),
    onSettled: () =>
      queryClient.invalidateQueries({ queryKey: getFindInvitableFriendsQueryKey(communityId) }),
  });
};
