import { Box, Stack, Typography } from "@mui/material";
import { useUserFriends } from "api/user/queries";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import React from "react";

interface IProps { }

export const Friends: React.FC<IProps> = () => {
  const { data } = useUserFriends();

  const items =
    data?.friends.map((f) => (
      <Stack direction={"row"} gap={1}>
        <img src={f.profilePicture} height={40} />
        <Typography>{f.firstName}</Typography>
      </Stack>
    )) ?? [];

  return (
    <Box width={"100%"}>
      <BoxWithTitle title={"meus amigos"}>
        <StripedGrid items={items} elementsPerRow={1} />
      </BoxWithTitle>
    </Box>
  );
};
