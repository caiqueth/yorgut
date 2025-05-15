import React from "react";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import { PicAndText } from "components/atoms/PicAndText";
import { useMeMemberships } from "api/me/queries";

interface ICommunitiesListProps { };
export const CommunitiesList: React.FC<ICommunitiesListProps> = () => {
  const { data: memberships } = useMeMemberships();

  const items =
    memberships?.communities.map((c) => (
      <PicAndText
        imageSrc={`/${c.profilePicture}`}
        url={`/community/${c.id}`}
        text={`${c.name} (${c.membershipCount ?? 0})`}
      />
    )) ?? [];

  return (
    <BoxWithTitle title={`minhas comunidades (${memberships?.membershipCount})`}>
      <StripedGrid items={items} elementsPerRow={3} />
    </BoxWithTitle>
  );
};
