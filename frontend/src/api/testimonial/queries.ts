import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { doDelete, doGet, doPost } from "api/http";
import { ICreateTestimonialPayload, ITestimonialsResponse } from "./types";

const basePath = `${__API_BASE_URL__}/testimonials`;

// Without page, returns every approved testimonial. With page (from 1), returns
// one page of them; the key keeps the base URL first so invalidating
// `${basePath}/user/${userId}` also refreshes every page.
export const useTestimonials = (userId: number, page?: number, pageSize?: number) =>
  useQuery<ITestimonialsResponse>({
    queryKey: page
      ? [`${basePath}/user/${userId}`, "page", page, pageSize]
      : [`${basePath}/user/${userId}`],
    queryFn: () =>
      doGet({
        queryKey: [
          page
            ? `${basePath}/user/${userId}?page=${page}&pageSize=${pageSize ?? 5}`
            : `${basePath}/user/${userId}`,
        ],
      }),
    enabled: !!userId,
    placeholderData: (previous) => previous,
  });

export const useCreateTestimonial = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: ICreateTestimonialPayload) =>
      doPost({ url: `${basePath}`, body: payload }),
    onSettled: (_data, _error, variables) => {
      queryClient.invalidateQueries({
        queryKey: [`${basePath}/user/${variables.userRecipientId}`],
      });
    },
  });
};

// userId is the recipient whose testimonials list should refresh.
export const useApproveTestimonial = (userId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (testimonialId: number) =>
      doPost({ url: `${basePath}/${testimonialId}/approve`, body: {} }),
    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: [`${basePath}/user/${userId}`],
      });
    },
  });
};

// userId is the recipient whose testimonials list should refresh.
export const useDeleteTestimonial = (userId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (testimonialId: number) =>
      doDelete({ url: `${basePath}/${testimonialId}` }),
    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: [`${basePath}/user/${userId}`],
      });
    },
  });
};
