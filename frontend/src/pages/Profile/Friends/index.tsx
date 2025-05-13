import { Box, Stack, Typography } from "@mui/material";
import { useUser, useUserFriends } from "api/user/queries";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import React from "react";
import { Link, useParams } from "react-router-dom";

interface IProps { }
export const UserFriendsList: React.FC<IProps> = () => {
  const params = useParams();
  const userId = parseInt(params.userId ?? "");
  const { data: friends } = useUserFriends(userId);
  const { data: user } = useUser(userId)

  const items = friends?.friends.map((f) => (
    <Link to={`/profile/${f.id}`}>
      <Stack direction={"row"} alignItems={"center"} gap={1}>
        <img src={`/${f.profilePicture}`} height={48}></img>
        <Typography>{f.firstName}</Typography>
      </Stack>
    </Link>
  )) ?? [];

  return (<Box>
    <BoxWithTitle title={`amigos de ${user?.firstName}`}>
      <StripedGrid items={items} elementsPerRow={1}></StripedGrid>
    </BoxWithTitle>
  </Box>);
};
