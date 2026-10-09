import { useMutation, useQueryClient } from "@tanstack/react-query";
import { loginWithGoogle } from "api/generated/endpoints/auth/auth";
import { LoginWithGoogleRequest } from "api/generated/model";

// Exchanges the Google credential for a Yorgut session token and stores it.
export const useLoginWithGoogle = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: LoginWithGoogleRequest) => loginWithGoogle(payload),
    onSuccess: ({ token }) => {
      localStorage.setItem("token", token);
      // Drop anything cached before login (e.g. the failed /me, or another
      // user's data) and refetch what's on screen with the new token.
      queryClient.resetQueries();
    },
  });
};
