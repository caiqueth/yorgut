import React from "react";

import { ProfileTable } from "components/ProfileTable";
import { useCommunityProfile } from "api/user/queries";

interface IProfileDetailsProps {
  communityId: number;
}
export const ProfileDetails: React.FC<IProfileDetailsProps> = ({
  communityId,
}) => {
  const { data: communityProfile } = useCommunityProfile(communityId);

  const data = [
    {
      label: "descrição:",
      value: communityProfile?.description,
    },
  ];

  return <ProfileTable data={data} />;
};
