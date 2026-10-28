import { useLogout } from "hooks/useLogout";

// The "Sair" link: ends the session instead of following the href.
export const LogoutLink: React.FC = () => {
  const logout = useLogout();

  const handleLogout = (e: React.MouseEvent) => {
    e.preventDefault();
    logout();
  };

  return <a href="/login" onClick={handleLogout}>Sair</a>;
};
