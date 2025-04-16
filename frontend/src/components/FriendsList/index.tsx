import React from "react";
import { useUserFriends, useUserStatistics } from "api/user/queries";
import { Button } from "components/Button";
import { isStripedRow } from "helpers/isStripedRow";
import { Link, useNavigate } from "react-router-dom";
import { palette } from "theme";

import { Box, Grid2 as Grid, Paper, Stack, Typography } from "@mui/material";
import { BoxWithTitle } from "components/BoxWithTitle";

interface IFriendsListProps {
  userId: number;
}
export const FriendsList: React.FC<IFriendsListProps> = ({ userId }) => {
  const userStatistics = useUserStatistics(userId);
  const { data: friends } = useUserFriends(userId);
  const navigate = useNavigate();

  const handleGoToFriends = () => {
    navigate("friends");
  };

  return (
    <BoxWithTitle title={`amigos (${userStatistics.totalFriends})`}>
      <Stack direction={"column"} alignItems={"center"}>
        <Grid container direction={"row"}>
          {friends?.friends?.map((friend, idx) => (
            <Grid
              size={3}
              key={idx}
              sx={{
                p: 1,
                backgroundColor: isStripedRow(idx, 4)
                  ? palette.darkBlue
                  : palette.mediumBlue,
              }}
            >
              <Stack
                alignItems={"center"}
                flexDirection={"column"}
                sx={{ textAlign: "center" }}
              >
                <Link to={`/profile/${friend.id}`}>
                  <img
                    height={64}
                    src={`/${friend.profilePicture}`}
                    alt="image not found"
                  />
                  <Typography
                    variant={"body2"}
                  >{`${friend.firstName} (${friend.friendCount})`}</Typography>
                </Link>
              </Stack>
            </Grid>
          ))}
        </Grid>
        <Box alignSelf={"center"}>
          <Button onClick={handleGoToFriends}>ver amigos</Button>
        </Box>
      </Stack>
    </BoxWithTitle>
  );
};
