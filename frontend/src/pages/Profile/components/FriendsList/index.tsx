import React from "react";
import { Button } from "components/atoms/Button";
import { useNavigate, useParams } from "react-router-dom";
import { Box, Stack } from "@mui/material";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import { PicAndText } from "components/atoms/PicAndText";
import { useMe, useMeFriends } from "api/me/queries";
import { useUser, useUserFriends } from "api/user/queries";

interface IFriendsListProps { }
export const FriendsList: React.FC<IFriendsListProps> = () => {
  const params = useParams();
  const userId = parseInt(params.userId ?? "");

  const { data: user } = useUser(userId);
  const { data: friends } = useUserFriends(userId);

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
