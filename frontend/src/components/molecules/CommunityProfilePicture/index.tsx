import React from "react";

import { ProfilePicture } from "components/atoms/ProfilePicture";
import { useCommunityProfile } from "api/community/queries";
import { Link } from "react-router-dom";

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
