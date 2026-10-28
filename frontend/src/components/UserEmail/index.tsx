import { useMe } from "api/me/queries";
import { useToken } from "hooks/useToken";
import { InlineSkeleton } from "components/InlineSkeleton";

// The signed-in user's e-mail, in bold. Logged out, useMe never asks for the
// user, so there's no email to wait for: nothing instead of a skeleton
// forever.
export const UserEmail: React.FC = () => {
  const { data: user } = useMe();
  const token = useToken();

  if (!token) return null;
  return <strong>{user ? user.email : <InlineSkeleton width={"20ch"} />}</strong>;
};
