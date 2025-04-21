import React from "react";
import { useUserMemberships, useUserStatistics } from "api/user/queries";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import { PicAndText } from "components/atoms/PicAndText";

interface ICommunitiesListProps {
  userId: number;
}
export const CommunitiesList: React.FC<ICommunitiesListProps> = ({
  userId,
}) => {
  const userStatistics = useUserStatistics(userId);
  const { data: memberships } = useUserMemberships(userId);

  const items =
    memberships?.communities.map((c) => (
      <PicAndText
        imageSrc={`/${c.profilePicture}`}
        url={`/community/${c.id}`}
        text={`${c.name} (${c.membershipCount ?? 0})`}
      />
    )) ?? [];

  return (
    <BoxWithTitle title={`comunidades (${userStatistics.totalCommunities})`}>
      <StripedGrid items={items} elementsPerRow={3} />
    </BoxWithTitle>
  );
};
