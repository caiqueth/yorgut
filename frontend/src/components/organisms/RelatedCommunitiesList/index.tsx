import React from "react";
import { useRelatedCommunities } from "api/user/queries";

import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { PicAndText } from "components/atoms/PicAndText";
import { StripedGrid } from "components/molecules/StripedGrid";

interface IRelatedCommunitiesListProps {
  communityId: number;
}
export const RelatedCommunitiesList: React.FC<IRelatedCommunitiesListProps> = ({
  communityId,
}) => {
  const { data: communities } = useRelatedCommunities(communityId);

  const items =
    communities?.communities.map((c) => (
      <PicAndText imageSrc={""} url={""} text={""} />
    )) ?? [];

  return (
    <BoxWithTitle
      title={`comunidades relacionadas (${communities?.memberCount})`}
    >
      <StripedGrid items={items} elementsPerRow={4} />
    </BoxWithTitle>
  );
};
