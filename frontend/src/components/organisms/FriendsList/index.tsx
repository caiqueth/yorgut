import React from "react";
import { useUserFriends, useUser } from "api/user/queries";
import { Button } from "components/atoms/Button";
import { isStripedRow } from "helpers/isStripedRow";
import { Link, useNavigate } from "react-router-dom";
import { palette } from "theme";

import { Box, Grid2 as Grid, Paper, Stack, Typography } from "@mui/material";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";

interface IFriendsListProps {
  userId: number;
}
export const FriendsList: React.FC<IFriendsListProps> = ({ userId }) => {
  const { data: user } = useUser(userId);
  const { data: friends } = useUserFriends(userId);
  const navigate = useNavigate();

  const handleGoToFriends = () => {
    navigate("friends");
  };

  const items =
    friends?.friends.map((f) => ({
      imageSrc: `/${f.profilePicture}`,
      url: `/profile/${f.id}`,
      text: `${f.firstName} (${f.friendCount})`,
    })) ?? [];

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
