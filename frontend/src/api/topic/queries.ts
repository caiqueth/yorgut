import { useQuery } from "@tanstack/react-query";
import { doGet } from "api/http";
import { ITopicRepliesResponse } from "./types";

const basePath = "/api/topic";

export const useTopicReplies = (topicId: number) =>
  useQuery<ITopicRepliesResponse>({
    queryKey: [`${basePath}/${topicId}/replies`],
    queryFn: ({ queryKey }) => doGet({ queryKey }),
    enabled: !!topicId,
  });
