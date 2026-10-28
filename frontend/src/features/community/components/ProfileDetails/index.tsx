import React from "react";
import { Skeleton } from "@mui/material";

import { FormattedText } from "components/FormattedText";
import { ProfileTable } from "components/ProfileTable";
import { useCommunity, useCommunityProfile } from "api/community/queries";
import { dayjsFromDate } from "helpers/formatting/date";
import { communityCategoryLabel } from "helpers/communityCategories";
import { Link } from "react-router-dom";

interface IProfileDetailsProps {
  communityId: number;
  // Continues the box above (ProfileHeader's joinedBelow): no top corners,
  // and no shadow over it. The top border stays, as the line between them.
  joinedAbove?: boolean;
}
export const ProfileDetails: React.FC<IProfileDetailsProps> = ({ communityId, joinedAbove }) => {
  const { data: community } = useCommunity(communityId);
  const { data: communityProfile } = useCommunityProfile(communityId);

  // The labels are fixed; only the values wait for the community.
  const loading = <Skeleton variant={"text"} width={"60%"} />;
  const data = [
    {
      label: "descrição:",
      value: !communityProfile
        ? loading
        : communityProfile.description && <FormattedText text={communityProfile.description} />,
    },
    {
      label: "categoria:",
      value: community ? communityCategoryLabel(community.category) : loading,
    },
    {
      label: "tipo:",
      value: community ? (community.private ? "🔒 privada (só membros veem o fórum)" : "pública") : loading,
    },
    {
      label: "dono(a):",
      value: community ? <Link to={`/profile/${community.ownerId}`}>{community.ownerFirstName}</Link> : loading,
    },
    {
      label: "criada em:",
      value: community ? dayjsFromDate(community.createdAt, "DD [de] MMMM [de] YYYY") : loading,
    },
  ];

  return (
    <ProfileTable
      data={data}
      sx={joinedAbove ? { borderTopLeftRadius: 0, borderTopRightRadius: 0, clipPath: "inset(0 -4px -4px -4px)" } : undefined}
    />
  );
};
