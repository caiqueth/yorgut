import React from "react";
import { useUserProfile } from "api/user/queries";

import { ProfilePicture } from "components/atoms/ProfilePicture";
import { Link } from "react-router-dom";

interface IUserProfilePictureProps { userId: number; };
export const UserProfilePicture: React.FC<IUserProfilePictureProps> = ({ userId }) => {
  const { data: user, isLoading } = useUserProfile(userId);

  return (
    <Link to={`/profile/${userId}`}>
      <ProfilePicture isLoading={isLoading} imageSrc={user?.profilePicture} />
    </Link>
  );
};
