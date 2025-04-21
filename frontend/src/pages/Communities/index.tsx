import { Box, Stack, Typography } from "@mui/material";
import { useLoggedUser, useUserMemberships } from "api/user/queries";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import React from "react";

interface IProps {}

export const Communities: React.FC<IProps> = () => {
  const { userId } = useLoggedUser();

  const { data } = useUserMemberships(userId);

  const items =
    data?.communities.map((c) => (
      <Stack direction={"row"} gap={1}>
        <img src={c.profilePicture} height={40} />
        <Typography>{c.name}</Typography>
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
