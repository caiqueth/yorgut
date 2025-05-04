import React from "react";
import { useUserMemberships, useUserStatistics } from "api/user/queries";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import { PicAndText } from "components/atoms/PicAndText";

interface ICommunitiesListProps { };
export const CommunitiesList: React.FC<ICommunitiesListProps> = () => {
  const userStatistics = useUserStatistics();
  const { data: memberships } = useUserMemberships();

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
