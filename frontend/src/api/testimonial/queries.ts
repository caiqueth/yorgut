import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  approveTestimonial,
  createTestimonial,
  deleteTestimonial,
  getFindUserTestimonialsQueryKey,
  useFindUserTestimonials,
} from "api/generated/endpoints/testimonial/testimonial";
import { CreateTestimonialRequest } from "api/generated/model";

// Without page, returns every approved testimonial. With page (from 1), returns
// one page of them. The key starts with the user's testimonials path either way,
// so invalidating getFindUserTestimonialsQueryKey(userId) refreshes every page.
export const useTestimonials = (userId: number, page?: number, pageSize?: number) =>
  useFindUserTestimonials(
    userId,
    page ? { page, pageSize: pageSize ?? 5 } : undefined,
    { query: { enabled: !!userId, placeholderData: (previous) => previous } },
  );

export const useCreateTestimonial = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateTestimonialRequest) => createTestimonial(payload),
    onSettled: (_data, _error, variables) => {
      queryClient.invalidateQueries({
        queryKey: getFindUserTestimonialsQueryKey(variables.userRecipientId),
      });
    },
  });
};

// userId is the recipient whose testimonials list should refresh.
export const useApproveTestimonial = (userId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (testimonialId: number) => approveTestimonial(testimonialId),
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: getFindUserTestimonialsQueryKey(userId) });
    },
  });
};

// userId is the recipient whose testimonials list should refresh.
export const useDeleteTestimonial = (userId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (testimonialId: number) => deleteTestimonial(testimonialId),
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: getFindUserTestimonialsQueryKey(userId) });
    },
  });
};
