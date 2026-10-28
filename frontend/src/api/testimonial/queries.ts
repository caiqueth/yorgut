import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  approveTestimonial,
  createTestimonial,
  deleteTestimonial,
  getFindUserTestimonialsQueryKey,
  useFindUserTestimonials,
} from "api/generated/endpoints/testimonial/testimonial";
import { CreateTestimonialRequest } from "api/generated/model";
import { usePreviewParams } from "contexts/ProfilePreviewContext";

// One page (from 1) of the approved testimonials. Every page's key starts with
// the user's testimonials path, so invalidating
// getFindUserTestimonialsQueryKey(userId) refreshes them all.
// enabled false skips the request, e.g. for testimonials kept to friends.
export const useTestimonials = (userId: number, page: number, pageSize: number, enabled = true) => {
  const preview = usePreviewParams();
  return useFindUserTestimonials(
    userId,
    { page, pageSize, ...preview },
    { query: { enabled: enabled && !!userId, placeholderData: (previous) => previous } },
  );
};

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
