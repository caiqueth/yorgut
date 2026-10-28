import { googleLogout } from "@react-oauth/google";
import { useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";

// Ends the session: forgets the token and everything cached for this user,
// then goes to the login page.
export const useLogout = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  return () => {
    googleLogout();
    localStorage.removeItem("token");
    queryClient.clear();
    navigate("/login");
  };
};
