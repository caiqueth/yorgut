import React from "react";
import { useRelatedCommunities } from "api/user/queries";

import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { PicAndText } from "components/atoms/PicAndText";
import { StripedGrid } from "components/molecules/StripedGrid";
import { Box, Stack } from "@mui/material";
import { Button } from "components/atoms/Button";
import { useNavigate } from "react-router-dom";

interface IRelatedCommunitiesListProps {
  communityId: number;
}
export const RelatedCommunitiesList: React.FC<IRelatedCommunitiesListProps> = ({
  communityId,
}) => {
  const navigate = useNavigate();
  const { data: communities } = useRelatedCommunities(communityId);

  const handleGoToRelated = () => {
    navigate("related");
  };

  const items =
    communities?.communities.map((c) => (
      <PicAndText imageSrc={""} url={""} text={""} />
    )) ?? [];

  return (
    <BoxWithTitle
      title={`comunidades relacionadas (${communities?.memberCount})`}
    >
      <Stack direction={"column"} alignItems={"center"}>
        <StripedGrid items={items} elementsPerRow={4} />
        <Box alignSelf={"center"}>
          <Button onClick={handleGoToRelated}>ver todas</Button>
        </Box>
      </Stack>
    </BoxWithTitle>
  );
};
