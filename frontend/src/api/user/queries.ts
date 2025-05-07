import { doGet } from "api/http";

import { useQuery } from "@tanstack/react-query";

import {
  IUserFriendsResponse,
  IUserMembershipsResponse,
  IUserPhotoAlbumsResponse,
  IUserProfileResponse,
  IUserResponse,
  IUsersResponse,
} from "./types";

const basePath = "/api/users";

export const useUsers = () =>
  useQuery<IUsersResponse>({
    queryKey: [`${basePath}`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useUser = (userId: number) =>
  useQuery<IUserResponse>({
    queryKey: [`${basePath}/${userId}`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useUserFriends = (userId: number) =>
  useQuery<IUserFriendsResponse>({
    queryKey: [`${basePath}/${userId}/friends`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useUserMemberships = (userId: number) =>
  useQuery<IUserMembershipsResponse>({
    queryKey: [`${basePath}/${userId}/memberships`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useUserPhotoAlbums = (userId: number) =>
  useQuery<IUserPhotoAlbumsResponse>({
    queryKey: [`${basePath}/${userId}/photoAlbums`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useUserProfile = (userId: number) =>
  useQuery<IUserProfileResponse>({
    queryKey: [`${basePath}/${userId}/profile`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useUserStatistics = (userId: number) => ({
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
    { name: "Aline" },
    { name: "Caique" },
    { name: "Carol" },
    { name: "Diego" },
    { name: "Lígia" },
    { name: "Luciana" },
    { name: "Renato" },
    { name: "Weslley" },
  ],
});

export const useTodaysLuck = () => `Boas notícias chegarão por e-mail.`;
