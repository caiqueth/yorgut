import React from "react";

import { ProfileTable } from "components/organisms/ProfileTable";
import { useCommunityProfile } from "api/user/queries";
import { useParams } from "react-router-dom";

interface IProfileDetailsProps {}
export const ProfileDetails: React.FC<IProfileDetailsProps> = ({}) => {
  const params = useParams();
  const communityId = parseInt(params.communityId ?? "");

  const { data: communityProfile } = useCommunityProfile(communityId);

  const data = [
    {
      label: "descrição:",
      value: communityProfile?.description,
    },
  ];

  return <ProfileTable data={data} />;
};
