import React from "react";
import { useCommunityProfile } from "api/user/queries";

import { ProfilePicture } from "components/atoms/ProfilePicture";

interface ICommunityProfilePictureProps {
  communityId: number;
}
export const CommunityProfilePicture: React.FC<
  ICommunityProfilePictureProps
> = ({ communityId }) => {
  const { data: community } = useCommunityProfile(communityId);

  if (!community?.profilePicture) return null;

  return <ProfilePicture imageSrc={`/${community.profilePicture}`} />;
};
