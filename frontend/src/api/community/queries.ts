import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { doGet, doPost, doPut } from "api/http";
import {
  ICommunitiesResponse,
  ICommunityHotTopicsResponse,
  ICommunityMembersResponse,
  ICommunityProfileResponse,
  ICommunityRelatedResponse,
  ICommunityResponse,
  ICommunityTopicsResponse,
  ICreateCommunityRequest,
  ICreateCommunityResponse,
  ISearchCommunitiesResponse,
  IUpdateCommunityPayload,
} from "./types";

const basePath = `${__API_BASE_URL__}/communities`;

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

// Everything that shows a community's name, category, description or picture.
const communityKeys = (communityId: number) => [
  `${basePath}/${communityId}`,
  `${basePath}/${communityId}/profile`,
  `${basePath}`,
  `${basePath}/search`,
  `${__API_BASE_URL__}/me/my-communities`,
  `${__API_BASE_URL__}/me/memberships`,
];

// Owner only; the backend answers 403 to anyone else.
export const useUpdateCommunity = (communityId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: IUpdateCommunityPayload) =>
      doPut({ url: `${basePath}/${communityId}`, body: payload }),
    onSettled: () => {
      for (const key of communityKeys(communityId)) {
        queryClient.invalidateQueries({ queryKey: [key] });
      }
    },
  });
};

// Owner only; the backend answers 403 to anyone else.
export const useUpdateCommunityPicture = (communityId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (picture: File) => {
      const form = new FormData();
      form.append("picture", picture);
      return doPut({ url: `${basePath}/${communityId}/picture`, body: form });
    },
    onSettled: () => {
      for (const key of communityKeys(communityId)) {
        queryClient.invalidateQueries({ queryKey: [key] });
      }
    },
  });
};

// Sent as multipart/form-data so the optional picture travels with the fields.
export const useCreateCommunity = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ name, category, picture }: ICreateCommunityRequest) => {
      const form = new FormData();
      form.append("name", name);
      form.append("category", category);
      if (picture) form.append("picture", picture);
      return doPost<ICreateCommunityResponse>({ url: `${basePath}`, body: form });
    },
    onSettled: () => {
      for (const key of [
        `${__API_BASE_URL__}/me/my-communities`,
        `${__API_BASE_URL__}/me/memberships`,
        `${basePath}`,
        `${basePath}/search`,
      ]) {
        queryClient.invalidateQueries({ queryKey: [key] });
      }
    },
  });
};

export const useDeleteCommunity = () => {
  throw new Error("Not implemented");
};

export const useCommunityHotTopics = (communityId: number) =>
  useQuery<ICommunityHotTopicsResponse>({
    queryKey: [`${basePath}/${communityId}/hot-topics`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useCommunityTopics = (communityId: number) =>
  useQuery<ICommunityTopicsResponse>({
    queryKey: [`${basePath}/${communityId}/topics`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

// Communities the logged user hasn't joined, filtered by name (empty = all).
// The base URL stays first in the key so joining/leaving can refresh every
// search at once by invalidating `${basePath}/search`.
export const useSearchCommunities = (name: string, page: number, pageSize: number) =>
  useQuery<ISearchCommunitiesResponse>({
    queryKey: [`${basePath}/search`, name, page, pageSize],
    queryFn: () =>
      doGet({
        queryKey: [
          `${basePath}/search?name=${encodeURIComponent(name)}&page=${page}&pageSize=${pageSize}`,
        ],
      }),
    placeholderData: (previous) => previous,
  });
