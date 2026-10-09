import { Box, Typography } from "@mui/material";
import { useMeMyCommunities } from "api/me/queries";
import { PicAndTextHorizontal } from "components/PicAndTextHorizontal";
import { BoxWithTitle } from "components/BoxWithTitle";
import { StripedGrid } from "components/StripedGrid";
import React from "react";

export const MyCommunities: React.FC = () => {
  const { data, isLoading } = useMeMyCommunities();

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
    <BoxWithTitle title="minhas comunidades">
      {!items.length ? (
        <Box sx={{ p: 1 }}>
          <Typography>Parece que você ainda não possui nenhuma comunidade.</Typography>
        </Box>
      ) : <StripedGrid items={items} elementsPerRow={1} isLoading={isLoading} />}
    </BoxWithTitle>
  )
}