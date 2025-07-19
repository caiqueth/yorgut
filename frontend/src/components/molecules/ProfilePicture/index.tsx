import React from "react";
import { useUserProfile } from "api/user/queries";

import { ProfilePicture } from "components/atoms/ProfilePicture";
import { useCommunityProfile } from "api/community/queries";
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

interface ICommunityProfilePictureProps {
  communityId: number;
}
export const CommunityProfilePicture: React.FC<
  ICommunityProfilePictureProps
> = ({ communityId }) => {
  const { data: community, isLoading } = useCommunityProfile(communityId);

  return (
    <Link to={`/community/${communityId}`}>
      <ProfilePicture isLoading={isLoading} imageSrc={community?.profilePicture} />
    </Link>
  );
};
