import React from "react";

import { BoxWithTitle } from "components/BoxWithTitle";
import { PicAndText, PicAndTextSkeleton } from "components/PicAndText";
import { StripedGrid } from "components/StripedGrid";
import { Stack } from "@mui/material";
import { useCommunityRelated } from "api/community/queries";

interface IRelatedCommunitiesListProps {
  communityId: number;
}
export const RelatedCommunitiesList: React.FC<IRelatedCommunitiesListProps> = ({
  communityId,
}) => {
  const { data: communities, isLoading, isFetching } = useCommunityRelated(communityId);

  const items =
    communities?.communities.map((c) => (
      <PicAndText key={c.id} imageSrc={c.profilePicture} url={`/community/${c.id}`} title={c.name} />
    )) ?? [];

  return (
    <BoxWithTitle title={`comunidades relacionadas`} showProgress={isFetching && !isLoading}>
      <Stack direction={"column"} alignItems={"center"}>
        <StripedGrid items={items} elementsPerRow={3} length={9} isLoading={isLoading} skeleton={<PicAndTextSkeleton />} />
      </Stack>
    </BoxWithTitle>
  );
};
