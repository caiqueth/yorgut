import { doGet } from "api/http";

import { useQuery } from "@tanstack/react-query";

import { ICommunityMembershipsResponse, IFriendsResponse, IUserProfileResponse, IUserResponse } from "./types";

export const useUser = (userId: number) =>
  useQuery<IUserResponse>({
    queryKey: [`/api/user/${userId}`, {}],
    queryFn: ({ queryKey }) => doGet<IUserResponse>({ queryKey }),
  });

  export const useUserProfile = (userId: number) =>
    useQuery<IUserProfileResponse>({
      queryKey: [`/api/user/${userId}/profile`, {}],
      queryFn: ({ queryKey }) => doGet<IUserProfileResponse>({ queryKey }),
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
    queryKey: [`/api/user/${userId}/friends`, {}],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useUserMemberships = (userId: number) =>
  useQuery<ICommunityMembershipsResponse>({
    queryKey: [`/api/user/${userId}/memberships`, {}],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useTodaysLuck = (userId: number) => "Boas notícias chegarão por e-mail";
