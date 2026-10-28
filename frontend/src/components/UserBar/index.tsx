import { Link } from "react-router-dom";
import { useToken } from "hooks/useToken";
import { LinkBar } from "components/LinkBar";
import { LogoutLink } from "components/LogoutLink";
import { UserEmail } from "components/UserEmail";

interface IUserBarProps {
  // False leaves plain gaps between the links (LinkBar).
  separators?: boolean;
}
export const UserBar: React.FC<IUserBarProps> = ({ separators }) => {
  // Logged out there's no e-mail item (UserEmail), and no "|" before it.
  const token = useToken();

  return (
    <LinkBar
      separators={separators}
      justifyContent={{ xs: "start", sm: "end" }}
      items={[
        ...(token ? [<UserEmail />] : []),
        <Link to="/my-account">Minha conta</Link>,
        <Link to="/help">Ajuda</Link>,
        <LogoutLink />,
      ]}
    />
  );
};
