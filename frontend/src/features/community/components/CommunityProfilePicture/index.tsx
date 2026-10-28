import React from "react";

import { ProfilePicture } from "components/ProfilePicture";
import { useCommunityProfile } from "api/community/queries";

interface ICommunityProfilePictureProps {
  communityId: number;
}
export const CommunityProfilePicture: React.FC<
  ICommunityProfilePictureProps
> = ({ communityId }) => {
  const { data: community, isLoading } = useCommunityProfile(communityId);

  return (
    <ProfilePicture
      isLoading={isLoading}
      imageSrc={community?.profilePicture}
      to={`/community/${communityId}`}
    />
  );
};
