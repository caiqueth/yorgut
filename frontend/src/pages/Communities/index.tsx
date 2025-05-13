import { Box, Stack, Typography } from "@mui/material";
import { useMeMemberships } from "api/me/queries";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import React from "react";
import { Link } from "react-router-dom";

interface IProps { }

export const Communities: React.FC<IProps> = () => {
  const { data } = useMeMemberships();

  const items =
    data?.communities.map((c) => (
      <Stack direction={"row"} gap={2} alignItems={"center"}>
        <img src={c.profilePicture} width={64} height={64} />
        <Link to={`/community/${c.id}`}>
          <Typography>{c.name}</Typography>
          <Typography>{c.membershipCount} {c.membershipCount > 1 ? "membros" : "membro"}</Typography>
        </Link>
      </Stack>
    )) ?? [];

  return (
    <Box width={"100%"}>
      <BoxWithTitle title={"comunidades que participo"}>
        <StripedGrid items={items} elementsPerRow={1} />
      </BoxWithTitle>
    </Box>
  );
};
