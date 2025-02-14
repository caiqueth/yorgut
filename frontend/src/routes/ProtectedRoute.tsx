import { Disconnected } from "pages/Disconnected";

import { Backdrop, CircularProgress } from "@mui/material";

const useAuth = () => ({
  isLoading: false,
  isAuthenticated: true,
});

interface IProtectedRouteProps {
  children: React.ReactNode;
}
export const ProtectedRoute: React.FC<IProtectedRouteProps> = ({ children }) => {
  const auth = useAuth();

  if (auth.isLoading)
    return (
      <Backdrop open>
        <CircularProgress />
      </Backdrop>
    );

  if (!auth.isAuthenticated) {
    return <Disconnected />;
  }

  return <>{children}</>;
};
