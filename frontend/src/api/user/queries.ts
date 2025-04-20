import { doGet, doPost } from "api/http";

import { useMutation, useQuery } from "@tanstack/react-query";

import {
  ICommunityMembershipsResponse,
  ICommunityMembersResponse,
  ICommunityProfileResponse,
  ICommunityResponse,
  IFriendsResponse,
  IPendingFriendshipsResponse,
  IRelatedCommunitiesResponse,
  IUserProfileResponse,
  IUserResponse,
} from "./types";

export const useLoggedUser = () => {
  return {
    userId: 7,
  };
};

export const useUser = (userId: number) =>
  useQuery<IUserResponse>({
    queryKey: [`/api/user/${userId}`],
    queryFn: ({ queryKey }) => doGet<IUserResponse>({ queryKey }),
  });

export const useUserProfile = (userId: number) =>
  useQuery<IUserProfileResponse>({
    queryKey: [`/api/user/${userId}/profile`],
    queryFn: ({ queryKey }) => doGet<IUserProfileResponse>({ queryKey }),
  });

export const useUserPendingFriendships = (userId: number) =>
  useQuery<IPendingFriendshipsResponse>({
    queryKey: [`/api/user/${userId}/friendships/pending`],
    queryFn: ({ queryKey }) => doGet<IPendingFriendshipsResponse>({ queryKey }),
  });

export const useUserStatistics = (userId: number) => ({
  userId: userId,
  totalProfileViews: 3558,
  lastWeekProfileViews: 45,
  totalFriends: 111,
  totalCommunities: 123,
  undirectConnections: 34537335,
  totalScraps: 123,
  totalPhotos: 12,
  totalVideos: 1,
  totalFans: 8,
  trustworthy: 1,
  cool: 1,
  beautiful: 1,
  profileLastVisitors: [
    { name: "Aline", birthDay: new Date(), friends: 123 },
    { name: "Caique", birthDay: new Date(), friends: 123 },
    { name: "Carol", birthDay: new Date(), friends: 123 },
    { name: "Diego", birthDay: new Date(), friends: 123 },
    { name: "Lígia", birthDay: new Date(), friends: 123 },
    { name: "Luciana", birthDay: new Date(), friends: 123 },
    { name: "Renato", birthDay: new Date(), friends: 123 },
    { name: "Weslley", birthDay: new Date(), friends: 123 },
  ],
});

export const useUserFriends = (userId: number) =>
  useQuery<IFriendsResponse>({
    queryKey: [`/api/user/${userId}/friends`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useUserMemberships = (userId: number) =>
  useQuery<ICommunityMembershipsResponse>({
    queryKey: [`/api/user/${userId}/memberships`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useJoinCommunity = () => {
  return useMutation({
    mutationFn: ({
      userId,
      communityId,
    }: {
      userId: number;
      communityId: number;
    }) => doPost({ url: `/api/user/${userId}/join/${communityId}`, body: {} }),
  });
};

export const useLeaveCommunity = () => {
  return useMutation({
    mutationFn: ({
      userId,
      communityId,
    }: {
      userId: number;
      communityId: number;
    }) => doPost({ url: `/api/user/${userId}/leave/${communityId}`, body: {} }),
  });
};
export const useTodaysLuck = (userId: number) =>
  "Boas notícias chegarão por e-mail";

export const useCommunity = (communityId: number) =>
  useQuery<ICommunityResponse>({
    queryKey: [`/api/community/${communityId}`],
    queryFn: ({ queryKey }) => doGet<ICommunityResponse>({ queryKey }),
  });

export const useRelatedCommunities = (communityId: number) =>
  useQuery<IRelatedCommunitiesResponse>({
    queryKey: [`/api/community/${communityId}/related`],
    queryFn: ({ queryKey }) => doGet<IRelatedCommunitiesResponse>({ queryKey }),
  });

export const useCommunityProfile = (communityId: number) =>
  useQuery<ICommunityProfileResponse>({
    queryKey: [`/api/community/${communityId}/profile`],
    queryFn: ({ queryKey }) => doGet<ICommunityProfileResponse>({ queryKey }),
  });

export const useCommunityMembers = (communityId: number) =>
  useQuery<ICommunityMembersResponse>({
    queryKey: [`/api/community/${communityId}/members`],
    queryFn: ({ queryKey }) => doGet<ICommunityMembersResponse>({ queryKey }),
  });
