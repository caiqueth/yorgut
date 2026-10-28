import React from "react";
import { useUserProfile } from "api/user/queries";

import { ProfilePicture } from "components/ProfilePicture";

interface IUserProfilePictureProps {
  // Undefined while the caller is still loading it (the home page waits for /me).
  userId?: number;
}
export const UserProfilePicture: React.FC<IUserProfilePictureProps> = ({ userId }) => {
  const { data: user, isLoading } = useUserProfile(userId ?? 0);
  return (
    <ProfilePicture
      isLoading={!userId || isLoading}
      imageSrc={user?.profilePicture}
      to={userId ? `/profile/${userId}` : undefined}
    />
  );
};
