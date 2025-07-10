import { doGet, doPost } from "api/http";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  IAmIFriendResponse,
  IAmIMemberResponse,
  IMeFriendsPendingResponse,
  IMeFriendsResponse,
  IMeMembershipsResponse,
  IMePhotoAlbumsResponse,
  IMeProfileResponse,
  IMeResponse,
  IMeUpdateResponse,
  WelcomeResponse,
} from "./types";

const basePath = "/api/me";

export const useMe = () =>
  useQuery<IMeResponse>({
    queryKey: [`${basePath}`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useUpdateMe = () =>
  useQuery<IMeUpdateResponse>({
    queryKey: [`${basePath}`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

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
        queryKey: [`/api/community/${variables.communityId}/members`],
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
        queryKey: [`/api/community/${variables.communityId}/members`],
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

export const useMePhotoAlbums = () =>
  useQuery<IMePhotoAlbumsResponse>({
    queryKey: [`${basePath}/photoAlbums`],
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
