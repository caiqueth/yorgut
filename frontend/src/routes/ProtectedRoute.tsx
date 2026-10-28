import { useToken } from "hooks/useToken";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

interface IProtectedRouteProps {
  children: React.ReactNode;
}
export const ProtectedRoute: React.FC<IProtectedRouteProps> = ({
  children,
}) => {
  const navigate = useNavigate();
  const token = useToken();

  useEffect(() => {
    if (!token) {
      navigate("/login")
    }
  }, [token])

  return <>{children}</>;
};
