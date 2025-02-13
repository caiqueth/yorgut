import React from "react";
import { useUserFriends, useUserStatistics } from "api/user/queries";
import { palette } from "theme";

import { Box, Grid2 as Grid, Paper, Stack, Typography } from "@mui/material";

interface IFriendsListProps {
  userId: number;
}
export const FriendsList: React.FC<IFriendsListProps> = ({ userId }) => {
  const userStatistics = useUserStatistics(userId);
  const friends = useUserFriends(userId);

  return (
    <Paper>
      <Box sx={{ textAlign: "center", width: "100%", backgroundColor: palette.lightBlue }}>
        <Typography variant={"subtitle2"} fontWeight={"bold"}>
          {`amigos (${userStatistics.totalFriends})`}
        </Typography>
      </Box>
      <Grid container direction={"row"} spacing={3} sx={{ p: 1 }}>
        {friends.map((friend, idx) => {
          return (
            <Grid size={3} key={idx}>
              <Stack alignItems={"center"} flexDirection={"column"} sx={{ textAlign: "center" }}>
                <img height={64} src={friend.profilePicture} alt="image not found" />
                <Typography variant={"body2"}>{`${friend.name} (${friend.friends})`}</Typography>
              </Stack>
            </Grid>
          );
        })}
      </Grid>
    </Paper>
  );
};
