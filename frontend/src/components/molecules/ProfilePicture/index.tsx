import React from "react";
import { useUserProfile } from "api/user/queries";

import { ProfilePicture } from "components/atoms/ProfilePicture";
import { useCommunityProfile } from "api/community/queries";
import { Link } from "react-router-dom";

interface IUserProfilePictureProps { userId: number; };
export const UserProfilePicture: React.FC<IUserProfilePictureProps> = ({ userId }) => {
  const { data: user } = useUserProfile(userId);

  if (!user) return null;

  return <ProfilePicture imageSrc={`/${user.profilePicture}`} />;
};

interface ICommunityProfilePictureProps {
  communityId: number;
}
export const CommunityProfilePicture: React.FC<
  ICommunityProfilePictureProps
> = ({ communityId }) => {
  const { data: community } = useCommunityProfile(communityId);

  if (!community?.profilePicture) return null;

  return <Link to={`/community/${communityId}`}><ProfilePicture imageSrc={`/${community.profilePicture}`} /></Link>;
};
