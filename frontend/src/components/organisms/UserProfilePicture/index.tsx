import React from "react";
import { useUserProfile } from "api/user/queries";

import { ProfilePicture } from "components/atoms/ProfilePicture";

interface IUserProfilePictureProps { };
export const UserProfilePicture: React.FC<IUserProfilePictureProps> = ({

}) => {
  const { data: user } = useUserProfile();

  if (!user?.profilePicture) return null;

  return <ProfilePicture imageSrc={`/${user.profilePicture}`} />;
};
