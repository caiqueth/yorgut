import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  getFindCommunityMembersQueryKey,
  getSearchCommunitiesQueryKey,
} from "api/generated/endpoints/community/community";
import {
  acceptCommunityInvitation,
  cancelJoinRequest,
  declineCommunityInvitation,
  deleteMe,
  getAmIMemberOfQueryKey,
  getFindMyCommunityInvitationsQueryKey,
  getFindMeQueryKey,
  getFindMyMembershipsQueryKey,
  getFindMyProfileQueryKey,
  getFindMyWelcomeQueryKey,
  joinCommunity,
  leaveCommunity,
  updateMe,
  updateMyPrivacy,
  updateMyTheme,
  updateMyProfilePicture,
  useAmIFriendOf,
  useAmIMemberOf,
  useFindMe,
  useFindMyBirthdays,
  useFindMyCommunities,
  useFindMyCommunityInvitations,
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
import { FindMeResponse, UpdateMeRequest, UpdateMyPrivacyRequest, UpdateMyThemeRequestTheme } from "api/generated/model";
import { useToken } from "hooks/useToken";

// Skipped without a session token: logged out (login and about pages), the
// request could only answer 401, which redirects to the login page.
export const useMe = () => {
  const token = useToken();
  return useFindMe({ query: { staleTime: Infinity, enabled: !!token } });
};

// userId is the logged user's id, used to refresh their public profile queries.
export const useUpdateMe = (userId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    meta: { showsErrorInline: true },
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
    meta: { showsErrorInline: true },
    mutationFn: (picture: File) => updateMyProfilePicture({ picture }),
    onSettled: () => {
      for (const queryKey of [getFindMyProfileQueryKey(), getFindUserProfileQueryKey(userId)]) {
        queryClient.invalidateQueries({ queryKey });
      }
    },
  });
};

// One page of the logged user's friends, most recently seen first.
export const useMeFriends = (page: number, pageSize: number) =>
  useFindMyFriends({ page, pageSize }, { query: { placeholderData: (previous) => previous } });

export const useMeFriendsPending = () => useFindMyPendingFriendships();

// Friends with a birthday in the next 30 days, soonest first.
export const useMeBirthdays = () => useFindMyBirthdays();

// Joining, leaving or cancelling a join request change the same queries.
const invalidateMembership = (
  queryClient: ReturnType<typeof useQueryClient>,
  communityId: number,
) => {
  queryClient.invalidateQueries({ queryKey: getAmIMemberOfQueryKey(communityId) });
  queryClient.invalidateQueries({ queryKey: getFindCommunityMembersQueryKey(communityId) });
  queryClient.invalidateQueries({ queryKey: getFindMyMembershipsQueryKey() });
  queryClient.invalidateQueries({ queryKey: getSearchCommunitiesQueryKey() });
  // The welcome box counts communities for its next steps.
  queryClient.invalidateQueries({ queryKey: getFindMyWelcomeQueryKey() });
};

// For a private community the answer is { pending: true }: the owner has to
// accept the request first.
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

export const useCancelJoinRequest = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ communityId }: { communityId: number }) => cancelJoinRequest(communityId),
    onSettled: (_data, _error, variables) =>
      invalidateMembership(queryClient, variables.communityId),
  });
};

// The invitations friends sent the logged user to join their communities.
export const useMyCommunityInvitations = () => useFindMyCommunityInvitations();

interface ICommunityInvitationVariables {
  invitationId: number;
  // The invitation's community, whose membership queries change.
  communityId: number;
}

// For a private community, an invitation from anyone but the owner answers
// { pending: true }: the owner has to accept the user first.
export const useAcceptCommunityInvitation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ invitationId }: ICommunityInvitationVariables) => acceptCommunityInvitation(invitationId),
    onSettled: (_data, _error, variables) => {
      queryClient.invalidateQueries({ queryKey: getFindMyCommunityInvitationsQueryKey() });
      invalidateMembership(queryClient, variables.communityId);
    },
  });
};

export const useDeclineCommunityInvitation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ invitationId }: ICommunityInvitationVariables) => declineCommunityInvitation(invitationId),
    onSettled: () =>
      queryClient.invalidateQueries({ queryKey: getFindMyCommunityInvitationsQueryKey() }),
  });
};

// One page of the communities the logged user is a member of, by name, or
// shuffled for the home box (the API keeps the same shuffle all day, so
// reloads don't reshuffle).
export const useMeMemberships = (page: number, pageSize: number, random = false) =>
  useFindMyMemberships(
    { page, pageSize, random },
    { query: { placeholderData: (previous) => previous } },
  );

// One page of the communities the logged user created, by name.
export const useMeMyCommunities = (page: number, pageSize: number) =>
  useFindMyCommunities({ page, pageSize }, { query: { placeholderData: (previous) => previous } });

export const useMeProfile = () => useFindMyProfile();

// Saves every privacy setting at once. Pending until everything reloads (the
// user, the welcome box's visitors, what the profile shows), so the
// checkboxes don't flip back.
export const useUpdateMyPrivacy = () => {
  const queryClient = useQueryClient();
  return useMutation({
    meta: { showsErrorInline: true },
    mutationFn: (settings: UpdateMyPrivacyRequest) => updateMyPrivacy(settings),
    onSettled: () => queryClient.invalidateQueries(),
  });
};

// Changes the look of the app. The cached user changes before the request, so
// the new theme shows at once; a failure puts the old one back.
export const useUpdateMyTheme = () => {
  const queryClient = useQueryClient();
  return useMutation({
    meta: { showsErrorInline: true },
    mutationFn: (theme: UpdateMyThemeRequestTheme) => updateMyTheme({ theme }),
    onMutate: async (theme) => {
      const queryKey = getFindMeQueryKey();
      await queryClient.cancelQueries({ queryKey });
      const previous = queryClient.getQueryData<FindMeResponse>(queryKey);
      if (previous) queryClient.setQueryData<FindMeResponse>(queryKey, { ...previous, theme });
      return { previous };
    },
    onError: (_error, _theme, context) => {
      if (context?.previous) queryClient.setQueryData(getFindMeQueryKey(), context.previous);
    },
    onSettled: () => queryClient.invalidateQueries({ queryKey: getFindMeQueryKey() }),
  });
};

export const useMeAmMember = (communityId: number) => useAmIMemberOf(communityId);

export const useMeAmFriend = (userId: number) => useAmIFriendOf(userId);

export const useMeWelcome = () => useFindMyWelcome();

// Deletes the logged user's account for good. The caller logs out on success
// (useLogout), which also clears the cache.
export const useDeleteMe = () =>
  useMutation({
    meta: { showsErrorInline: true },
    mutationFn: () => deleteMe(),
  });
