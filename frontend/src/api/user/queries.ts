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

const basePath = `${__API_BASE_URL__}/users`;

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
    enabled: !!userId,
  });
