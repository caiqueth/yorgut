import { Box, Stack, Typography } from "@mui/material";
import { useUser, useUserFriends } from "api/user/queries";
import { PicAndTextHorizontal } from "components/atoms/PicAndTextHorizontal";
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
    <PicAndTextHorizontal
      title={f.firstName}
      text={`${f.friendCount} amigos`}
      imageSrc={`/${f.profilePicture}`}
      imageSize={56}
      url={`/profile/${f.id}`}
    />
  )) ?? [];

  return (<Box>
    <BoxWithTitle title={`amigos de ${user?.firstName}`}>
      <StripedGrid items={items} elementsPerRow={1}></StripedGrid>
    </BoxWithTitle>
  </Box>);
};
