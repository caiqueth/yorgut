import React from "react";
import aline from "assets/users/aline.jpg";
import caique from "assets/users/caique.jpg";
import carol from "assets/users/carol.jpg";
import diego from "assets/users/diego.jpg";
import ligia from "assets/users/ligia.jpg";
import luciana from "assets/users/luciana.jpg";
import renato from "assets/users/renato.jpg";
import weslley from "assets/users/weslley.jpg";
import { palette } from "theme";

import { Box, Grid2 as Grid, Paper, Stack, Typography } from "@mui/material";

interface IFriendsListProps {
  userId: number;
}
export const FriendsList: React.FC<IFriendsListProps> = ({ userId }) => {
  const user = {
    friendsNumber: 132,
    communitiesNumber: 821,
  };

  const friends = [
    { name: "Aline", profilePicture: aline, friends: 123 },
    { name: "Caique", profilePicture: caique, friends: 123 },
    { name: "Carol", profilePicture: carol, friends: 123 },
    { name: "Diego", profilePicture: diego, friends: 123 },
    { name: "Lígia", profilePicture: ligia, friends: 123 },
    { name: "Luciana", profilePicture: luciana, friends: 123 },
    { name: "Renato", profilePicture: renato, friends: 123 },
    { name: "Weslley", profilePicture: weslley, friends: 123 },
  ];

  return (
    <Paper>
      <Box sx={{ textAlign: "center", width: "100%", backgroundColor: palette.lightBlue }}>
        <Typography variant={"subtitle2"} fontWeight={"bold"}>
          {`amigos (${user.friendsNumber})`}
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
