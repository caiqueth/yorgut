import { useQuery } from "@tanstack/react-query";
import { doGet } from "api/http";
import {
  ICommunitiesResponse,
  ICommunityHotTopicsResponse,
  ICommunityMembersResponse,
  ICommunityProfileResponse,
  ICommunityRelatedResponse,
  ICommunityResponse,
  ICommunityTopicsResponse,
} from "./types";

const basePath = "/api/community";

export const useCommunities = () =>
  useQuery<ICommunitiesResponse>({
    queryKey: [`${basePath}`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useCommunity = (communityId: number) =>
  useQuery<ICommunityResponse>({
    queryKey: [`${basePath}/${communityId}`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useCommunityMembers = (communityId: number) =>
  useQuery<ICommunityMembersResponse>({
    queryKey: [`${basePath}/${communityId}/members`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useCommunityProfile = (communityId: number) =>
  useQuery<ICommunityProfileResponse>({
    queryKey: [`${basePath}/${communityId}/profile`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useCommunityRelated = (communityId: number) =>
  useQuery<ICommunityRelatedResponse>({
    queryKey: [`${basePath}/${communityId}/related`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useCreateCommunityProfile = () => {
  throw new Error("Not implemented");
};

export const useUpdateCommunity = () => {
  throw new Error("Not implemented");
};

export const useCreateCommunity = () => {
  throw new Error("Not implemented");
};

export const useDeleteCommunity = () => {
  throw new Error("Not implemented");
};

export const useCommunityHotTopics = (communityId: number) =>
  useQuery<ICommunityHotTopicsResponse>({
    queryKey: [`${basePath}/${communityId}/hot_topics`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useCommunityTopics = (communityId: number) =>
  useQuery<ICommunityTopicsResponse>({
    queryKey: [`${basePath}/${communityId}/topics`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });
