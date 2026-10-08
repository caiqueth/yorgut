import { doGet, doPost, doPut } from "api/http";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  IAmIFriendResponse,
  IAmIMemberResponse,
  IMeFriendsPendingResponse,
  IMeFriendsResponse,
  IMeMembershipsResponse,
  IMeProfileResponse,
  IMeResponse,
  IUpdateMePayload,
  WelcomeResponse,
} from "./types";

const basePath = `${__API_BASE_URL__}/me`;

export const useMe = () =>
  useQuery<IMeResponse>({
    queryKey: [`${basePath}`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
    staleTime: Infinity,
  })


// userId is the logged user's id, used to refresh their public profile queries.
export const useUpdateMe = (userId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: IUpdateMePayload) =>
      doPut({ url: `${basePath}`, body: payload }),
    onSettled: () => {
      for (const key of [
        `${basePath}`,
        `${basePath}/profile`,
        `${basePath}/welcome`,
        `${__API_BASE_URL__}/users/${userId}`,
        `${__API_BASE_URL__}/users/${userId}/profile`,
      ]) {
        queryClient.invalidateQueries({ queryKey: [key] });
      }
    },
  });
};

// userId is the logged user's id, used to refresh their public profile queries.
export const useUpdateProfilePicture = (userId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (picture: File) => {
      const form = new FormData();
      form.append("picture", picture);
      return doPut({ url: `${basePath}/profile-picture`, body: form });
    },
    onSettled: () => {
      for (const key of [
        `${basePath}/profile`,
        `${__API_BASE_URL__}/users/${userId}/profile`,
      ]) {
        queryClient.invalidateQueries({ queryKey: [key] });
      }
    },
  });
};

export const useMeFriends = () =>
  useQuery<IMeFriendsResponse>({
    queryKey: [`${basePath}/friends`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useMeFriendsPending = () =>
  useQuery<IMeFriendsPendingResponse>({
    queryKey: [`${basePath}/friends/pending`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useJoinCommunity = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ communityId }: { communityId: number }) =>
      doPost({ url: `${basePath}/join/${communityId}`, body: {} }),
    onSettled: (_data, _error, variables) => {
      queryClient.invalidateQueries({
        queryKey: [`${basePath}/memberships/${variables.communityId}`],
      });
      queryClient.invalidateQueries({
        queryKey: [`${__API_BASE_URL__}/communities/${variables.communityId}/members`],
      });
      queryClient.invalidateQueries({
        queryKey: [`${basePath}/memberships`],
      });
      queryClient.invalidateQueries({
        queryKey: [`${__API_BASE_URL__}/communities/search`],
      });
    },
  });
};

export const useLeaveCommunity = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ communityId }: { communityId: number }) =>
      doPost({ url: `${basePath}/leave/${communityId}`, body: {} }),
    onSettled: (_data, _error, variables) => {
      queryClient.invalidateQueries({
        queryKey: [`${basePath}/memberships/${variables.communityId}`],
      });
      queryClient.invalidateQueries({
        queryKey: [`${__API_BASE_URL__}/communities/${variables.communityId}/members`],
      });
      queryClient.invalidateQueries({
        queryKey: [`${basePath}/memberships`],
      });
      queryClient.invalidateQueries({
        queryKey: [`${__API_BASE_URL__}/communities/search`],
      });
    },
  });
};

export const useMeMemberships = () =>
  useQuery<IMeMembershipsResponse>({
    queryKey: [`${basePath}/memberships`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useMeMyCommunities = () =>
  useQuery<IMeMembershipsResponse>({
    queryKey: [`${basePath}/my-communities`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useMeProfile = () =>
  useQuery<IMeProfileResponse>({
    queryKey: [`${basePath}/profile`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useMeAmMember = (communityId: number) =>
  useQuery<IAmIMemberResponse>({
    queryKey: [`${basePath}/memberships/${communityId}`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useMeAmFriend = (userId: number) =>
  useQuery<IAmIFriendResponse>({
    queryKey: [`${basePath}/friends/${userId}`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useMeWelcome = () =>
  useQuery<WelcomeResponse>({
    queryKey: [`${basePath}/welcome`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });
