import { Box, Stack, Typography } from "@mui/material";
import { useMeMemberships } from "api/me/queries";
import { PicAndTextHorizontal } from "components/atoms/PicAndTextHorizontal";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import React from "react";
import { Link } from "react-router-dom";

interface IProps { }

export const Communities: React.FC<IProps> = () => {
  const { data } = useMeMemberships();

  const items =
    data?.communities.map((c) => (
      <PicAndTextHorizontal
        title={c.name}
        text={`${c.membershipCount} membros`}
        url={`/community/${c.id}`}
        imageSrc={c.profilePicture}
        imageSize={64}
      />
    )) ?? [];

  return (
    <Box width={"100%"}>
      <BoxWithTitle title={"comunidades que participo"}>
        <StripedGrid items={items} elementsPerRow={1} />
      </BoxWithTitle>
    </Box>
  );
};
