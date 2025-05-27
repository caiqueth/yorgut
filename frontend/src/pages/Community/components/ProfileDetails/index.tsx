import React from "react";

import { ProfileTable } from "components/organisms/ProfileTable";
import { useCommunity, useCommunityProfile } from "api/community/queries";
import { dayjsFromDate } from "helpers/formatting/date";
import { Link } from "react-router-dom";

interface IProfileDetailsProps {
  communityId: number;
}
export const ProfileDetails: React.FC<IProfileDetailsProps> = ({ communityId }) => {
  const { data: community } = useCommunity(communityId);
  const { data: communityProfile } = useCommunityProfile(communityId);

  if (!communityProfile || !community) return null;

  const data = [
    {
      label: "descrição:",
      value: communityProfile?.description,
    },
    {
      label: "dono(a):",
      value: <Link to={`/profile/${community.ownerId}`}>{community?.ownerFirstName}</Link>,
    },
    {
      label: "criada em:",
      value: dayjsFromDate(community.createdAt, "DD [de] MMMM [de] YYYY"),
    },
  ];

  return <ProfileTable data={data} />;
};
