import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  getFindCommunityMembersQueryKey,
  getSearchCommunitiesQueryKey,
} from "api/generated/endpoints/community/community";
import {
  getAmIMemberOfQueryKey,
  getFindMeQueryKey,
  getFindMyMembershipsQueryKey,
  getFindMyProfileQueryKey,
  getFindMyWelcomeQueryKey,
  joinCommunity,
  leaveCommunity,
  updateMe,
  updateMyProfilePicture,
  useAmIFriendOf,
  useAmIMemberOf,
  useFindMe,
  useFindMyCommunities,
  useFindMyFriends,
  useFindMyMemberships,
  useFindMyPendingFriendships,
  useFindMyProfile,
  useFindMyWelcome,
} from "api/generated/endpoints/me/me";
import {
  getFindUserProfileQueryKey,
  getFindUserQueryKey,
} from "api/generated/endpoints/users/users";
import { UpdateMeRequest } from "api/generated/model";

export const useMe = () => useFindMe({ query: { staleTime: Infinity } });

// userId is the logged user's id, used to refresh their public profile queries.
export const useUpdateMe = (userId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: UpdateMeRequest) => updateMe(payload),
    onSettled: () => {
      for (const queryKey of [
        getFindMeQueryKey(),
        getFindMyProfileQueryKey(),
        getFindMyWelcomeQueryKey(),
        getFindUserQueryKey(userId),
        getFindUserProfileQueryKey(userId),
      ]) {
        queryClient.invalidateQueries({ queryKey });
      }
    },
  });
};

// userId is the logged user's id, used to refresh their public profile queries.
export const useUpdateProfilePicture = (userId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (picture: File) => updateMyProfilePicture({ picture }),
    onSettled: () => {
      for (const queryKey of [getFindMyProfileQueryKey(), getFindUserProfileQueryKey(userId)]) {
        queryClient.invalidateQueries({ queryKey });
      }
    },
  });
};

export const useMeFriends = () => useFindMyFriends();

export const useMeFriendsPending = () => useFindMyPendingFriendships();

// Joining or leaving changes the same queries.
const invalidateMembership = (
  queryClient: ReturnType<typeof useQueryClient>,
  communityId: number,
) => {
  queryClient.invalidateQueries({ queryKey: getAmIMemberOfQueryKey(communityId) });
  queryClient.invalidateQueries({ queryKey: getFindCommunityMembersQueryKey(communityId) });
  queryClient.invalidateQueries({ queryKey: getFindMyMembershipsQueryKey() });
  queryClient.invalidateQueries({ queryKey: getSearchCommunitiesQueryKey() });
};

export const useJoinCommunity = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ communityId }: { communityId: number }) => joinCommunity(communityId),
    onSettled: (_data, _error, variables) =>
      invalidateMembership(queryClient, variables.communityId),
  });
};

export const useLeaveCommunity = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ communityId }: { communityId: number }) => leaveCommunity(communityId),
    onSettled: (_data, _error, variables) =>
      invalidateMembership(queryClient, variables.communityId),
  });
};

export const useMeMemberships = () => useFindMyMemberships();

export const useMeMyCommunities = () => useFindMyCommunities();

export const useMeProfile = () => useFindMyProfile();

export const useMeAmMember = (communityId: number) => useAmIMemberOf(communityId);

export const useMeAmFriend = (userId: number) => useAmIFriendOf(userId);

export const useMeWelcome = () => useFindMyWelcome();
