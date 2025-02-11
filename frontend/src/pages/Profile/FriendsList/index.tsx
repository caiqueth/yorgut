import React from "react";

import { Box, Grid2 as Grid, Typography } from "@mui/material";

interface IFriendsListProps {
  userId: number;
}
export const FriendsList: React.FC<IFriendsListProps> = ({ userId }) => {
  const friends = [
    { name: "Aline" },
    { name: "Carol" },
    { name: "Diego" },
    { name: "Lígia" },
    { name: "Luciana" },
    { name: "Renato" },
    { name: "Weslley" },
  ];

  return (
    <Grid container direction={"row"} gap={3}>
      {friends.map((friend, idx) => (
        <Grid sx={{ xm: 3 }} key={idx}>
          <Box>
            <image height={48} />
            <Typography>{friend.name}</Typography>
          </Box>
        </Grid>
      ))}
    </Grid>
  );
};
