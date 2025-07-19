import React from "react";

import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { PicAndText } from "components/atoms/PicAndText";
import { StripedGrid } from "components/molecules/StripedGrid";
import { Stack } from "@mui/material";
import { useCommunityRelated } from "api/community/queries";

interface IRelatedCommunitiesListProps {
  communityId: number;
}
export const RelatedCommunitiesList: React.FC<IRelatedCommunitiesListProps> = ({
  communityId,
}) => {
  const { data: communities, isLoading } = useCommunityRelated(communityId);

  const items =
    communities?.communities.map((c) => (
      <PicAndText imageSrc={c.profilePicture} url={`/community/${c.id}`} text={c.name} />
    )) ?? [];

  return (
    <BoxWithTitle title={`comunidades relacionadas`}>
      <Stack direction={"column"} alignItems={"center"}>
        <StripedGrid items={items} elementsPerRow={4} isLoading={isLoading} />
      </Stack>
    </BoxWithTitle>
  );
};
