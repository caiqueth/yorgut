import { useMutation, useQueryClient } from "@tanstack/react-query";
import { ILoginResponse, ILoginWithGooglePayload } from "./types";
import { doPost } from "api/http";

const basePath = `${__API_BASE_URL__}/auth`;

// Exchanges the Google credential for a Yorgut session token and stores it.
export const useLoginWithGoogle = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: ILoginWithGooglePayload) =>
      doPost<ILoginResponse>({ url: `${basePath}/google`, body: payload }),
    onSuccess: ({ token }) => {
      localStorage.setItem("token", token);
      // Drop anything cached before login (e.g. the failed /me, or another
      // user's data) and refetch what's on screen with the new token.
      queryClient.resetQueries();
    },
  });
};
