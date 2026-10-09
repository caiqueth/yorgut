import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  createCommunity,
  getFindCommunitiesQueryKey,
  getFindCommunityProfileQueryKey,
  getFindCommunityQueryKey,
  getSearchCommunitiesQueryKey,
  updateCommunity,
  updateCommunityPicture,
  useFindCommunities,
  useFindCommunity,
  useFindCommunityHotTopics,
  useFindCommunityMembers,
  useFindCommunityProfile,
  useFindCommunityTopics,
  useFindRelatedCommunities,
  useSearchCommunities as useSearchCommunitiesQuery,
} from "api/generated/endpoints/community/community";
import {
  getFindMyCommunitiesQueryKey,
  getFindMyMembershipsQueryKey,
} from "api/generated/endpoints/me/me";
import { CreateCommunityBody, UpdateCommunityRequest } from "api/generated/model";

export const useCommunities = () => useFindCommunities();

export const useCommunity = (communityId: number) => useFindCommunity(communityId);

export const useCommunityMembers = (communityId: number) =>
  useFindCommunityMembers(communityId);

export const useCommunityProfile = (communityId: number) =>
  useFindCommunityProfile(communityId);

export const useCommunityRelated = (communityId: number) =>
  useFindRelatedCommunities(communityId);

// Everything that shows a community's name, category, description or picture.
const communityKeys = (communityId: number) => [
  getFindCommunityQueryKey(communityId),
  getFindCommunityProfileQueryKey(communityId),
  getFindCommunitiesQueryKey(),
  getSearchCommunitiesQueryKey(),
  getFindMyCommunitiesQueryKey(),
  getFindMyMembershipsQueryKey(),
];

// Owner only; the backend answers 403 to anyone else.
export const useUpdateCommunity = (communityId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
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
    mutationFn: (payload: CreateCommunityBody) => createCommunity(payload),
    onSettled: () => {
      for (const queryKey of [
        getFindMyCommunitiesQueryKey(),
        getFindMyMembershipsQueryKey(),
        getFindCommunitiesQueryKey(),
        getSearchCommunitiesQueryKey(),
      ]) {
        queryClient.invalidateQueries({ queryKey });
      }
    },
  });
};

export const useDeleteCommunity = () => {
  throw new Error("Not implemented");
};

export const useCommunityHotTopics = (communityId: number) =>
  useFindCommunityHotTopics(communityId);

export const useCommunityTopics = (communityId: number) =>
  useFindCommunityTopics(communityId);

// Communities the logged user hasn't joined, filtered by name (empty = all).
// The key starts with the search path, so joining/leaving can refresh every
// search at once by invalidating getSearchCommunitiesQueryKey().
export const useSearchCommunities = (name: string, page: number, pageSize: number) =>
  useSearchCommunitiesQuery(
    { name, page, pageSize },
    { query: { placeholderData: (previous) => previous } },
  );
