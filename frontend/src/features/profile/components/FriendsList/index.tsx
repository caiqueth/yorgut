import React from "react";
import { Button } from "components/Button";
import { useNavigate } from "react-router-dom";
import { Box, Stack } from "@mui/material";
import { BoxWithTitle } from "components/BoxWithTitle";
import { StripedGrid } from "components/StripedGrid";
import { PicAndText } from "components/PicAndText";
import { useUser, useUserFriends } from "api/user/queries";

interface IFriendsListProps {
  userId: number;
}
export const FriendsList: React.FC<IFriendsListProps> = ({ userId }) => {
  const { data: user } = useUser(userId);
  const { data: friends, isLoading } = useUserFriends(userId);

  const navigate = useNavigate();

  const handleGoToFriends = () => {
    navigate("friends");
  };

  const items =
    friends?.friends.map((f) => (
      <PicAndText
        imageSrc={f.profilePicture}
        url={`/profile/${f.id}`}
        text={`${f.firstName} (${f.friendCount})`}
      />
    )) ?? [];

  return (
    <BoxWithTitle title={`amigos (${user?.friendCount})`}>
      <Stack direction={"column"} alignItems={"center"}>
        <StripedGrid items={items} elementsPerRow={4} isLoading={isLoading} />
        <Box alignSelf={"center"}>
          <Button onClick={handleGoToFriends}>ver amigos</Button>
        </Box>
      </Stack>
    </BoxWithTitle>
  );
};
