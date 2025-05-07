import React from "react";
import { useUserProfile } from "api/user/queries";

import { ProfilePicture } from "components/atoms/ProfilePicture";

interface IUserProfilePictureProps {
  userId: number;
};
export const UserProfilePicture: React.FC<IUserProfilePictureProps> = ({
  userId
}) => {
  const { data: user } = useUserProfile(userId);

  if (!user?.profilePicture) return null;

  return <ProfilePicture imageSrc={`/${user.profilePicture}`} />;
};
