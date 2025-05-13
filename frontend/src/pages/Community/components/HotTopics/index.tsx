import { Box } from "@mui/material";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import React from "react";

interface IProps {
  communityId: number;
}
export const HotTopics: React.FC<IProps> = ({ communityId }) => {
  return <BoxWithTitle title={"tópicos"}>
    <StripedGrid items={[]} elementsPerRow={1} />
  </BoxWithTitle>
}