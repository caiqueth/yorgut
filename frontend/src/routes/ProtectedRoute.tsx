import { useToken } from "hooks/useToken";
import { useNavigate } from "react-router-dom";

interface IProtectedRouteProps {
  children: React.ReactNode;
}
export const ProtectedRoute: React.FC<IProtectedRouteProps> = ({
  children,
}) => {
  const navigate = useNavigate();
  const token = useToken();

  if (!token) {
    navigate("/login")
  }

  // if (auth.isLoading)
  //   return (
  //     <Backdrop open>
  //       <CircularProgress />
  //     </Backdrop>
  //   );

  // if (!auth.isAuthenticated) {
  //   return <Disconnected />;
  // }

  return children;
};
