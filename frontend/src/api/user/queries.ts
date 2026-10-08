import { doGet } from "api/http";

import { useQuery } from "@tanstack/react-query";

import {
  ISearchUsersResponse,
  IUserFriendsResponse,
  IUserMembershipsResponse,
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

export const useUserProfile = (userId: number) =>
  useQuery<IUserProfileResponse>({
    queryKey: [`${basePath}/${userId}/profile`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
    enabled: !!userId,
  });

export const useUserSearchUsers = (name: string, page: number, pageSize: number) =>
  useQuery<ISearchUsersResponse>({
    queryKey: [`${basePath}/search-users?name=${name}&page=${page}&pageSize=${pageSize}`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

