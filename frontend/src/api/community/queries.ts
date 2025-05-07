import { useQuery } from "@tanstack/react-query";
import { doGet } from "api/http";

export const useTodaysLuck = () => "Boas notícias chegarão por e-mail";

export const useCommunity = (communityId: number) =>
  useQuery({
    queryKey: [`/api/community/${communityId}`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useRelatedCommunities = (communityId: number) =>
  useQuery({
    queryKey: [`/api/community/${communityId}/related`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useCommunityProfile = (communityId: number) =>
  useQuery({
    queryKey: [`/api/community/${communityId}/profile`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });

export const useCommunityMembers = (communityId: number) =>
  useQuery({
    queryKey: [`/api/community/${communityId}/members`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
  });
