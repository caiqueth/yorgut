import React from "react";

import { useUserProfile } from "api/user/queries";
import { ProfileTable } from "components/ProfileTable";
import { profileFields } from "features/profile/profileFields";

interface IProfileDetailsProps {
  userId: number;
}
export const ProfileDetails: React.FC<IProfileDetailsProps> = ({ userId }) => {

  const { data: userProfile } = useUserProfile(userId);

  const data = profileFields.map(({ key, label }) => ({
    label,
    value: userProfile?.[key],
  }));

  return <ProfileTable data={data} />;
};
