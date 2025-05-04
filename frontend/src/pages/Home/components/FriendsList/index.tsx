import React from "react";
import { useUserFriends, useUser } from "api/user/queries";
import { Button } from "components/atoms/Button";
import { useNavigate } from "react-router-dom";
import { Box, Stack } from "@mui/material";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import { PicAndText } from "components/atoms/PicAndText";

interface IFriendsListProps { };

export const FriendsList: React.FC<IFriendsListProps> = () => {
  const { data: user } = useUser();
  const { data: friends } = useUserFriends();
  const navigate = useNavigate();

  const handleGoToFriends = () => {
    navigate("friends");
  };

  const items =
    friends?.friends.map((f) => (
      <PicAndText
        imageSrc={`/${f.profilePicture}`}
        url={`/profile/${f.id}`}
        text={`${f.firstName} (${f.friendCount})`}
      />
    )) ?? [];

  return (
    <BoxWithTitle title={`amigos (${user?.friendCount})`}>
      <Stack direction={"column"} alignItems={"center"}>
        <StripedGrid items={items} elementsPerRow={4} />
        <Box alignSelf={"center"}>
          <Button onClick={handleGoToFriends}>ver amigos</Button>
        </Box>
      </Stack>
    </BoxWithTitle>
  );
};
