import React from "react";
import { Skeleton } from "@mui/material";

import { useUserProfile } from "api/user/queries";
import { FormattedText } from "components/FormattedText";
import { ProfileTable } from "components/ProfileTable";
import { dayjsFromDate } from "helpers/formatting/date";
import { profileFields } from "features/profile/profileFields";

interface IProfileDetailsProps {
  userId: number;
  // Continues the box above (ProfileHeader's joinedBelow): no top corners,
  // and no shadow over it. The top border stays, as the line between them.
  joinedAbove?: boolean;
}
// The profile's fields. Fields the user keeps to friends are left out for
// everyone else (the API lists them in hiddenFields), instead of showing "-".
export const ProfileDetails: React.FC<IProfileDetailsProps> = ({ userId, joinedAbove }) => {
  const { data: userProfile } = useUserProfile(userId);

  // The labels are fixed; only the values wait for the profile.
  const loading = <Skeleton variant={"text"} width={"60%"} />;
  const hidden = userProfile?.hiddenFields ?? [];

  const fields = profileFields
    .filter(({ key }) => !hidden.includes(key))
    .map(({ key, label, formatted }) => {
      const value = userProfile?.[key];
      return {
        label,
        value: !userProfile ? loading : formatted && value ? <FormattedText text={value} /> : value,
      };
    });

  // The birthday goes after the relationship status; day and month only.
  const birthday = hidden.includes("birthDate")
    ? []
    : [{
      label: "aniversário:",
      value: !userProfile
        ? loading
        : userProfile.birthDate && dayjsFromDate(userProfile.birthDate, "D [de] MMMM"),
    }];
  const afterRelationship = fields.findIndex(({ label }) => label === "relacionamento:") + 1;
  const data = [...fields.slice(0, afterRelationship), ...birthday, ...fields.slice(afterRelationship)];

  return (
    <ProfileTable
      data={data}
      sx={joinedAbove ? { borderTopLeftRadius: 0, borderTopRightRadius: 0, clipPath: "inset(0 -4px -4px -4px)" } : undefined}
    />
  );
};
