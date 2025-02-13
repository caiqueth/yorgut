import React from "react";
import { useUserFriends } from "api/user/queries";
import { isStripedRow } from "helpers/isStripedRow";
import { palette } from "theme";

import { Box, Grid2 as Grid, Paper, Stack, Typography } from "@mui/material";

interface IBirthdaysProps {
  userId: number;
}
export const Birthdays: React.FC<IBirthdaysProps> = ({ userId }) => {
  const friends = useUserFriends(userId);
  return (
    <Paper>
      <Box sx={{ textAlign: "center", width: "100%", backgroundColor: palette.lightBlue }}>
        <Typography variant={"subtitle2"} fontWeight={"bold"}>
          {`próximos aniversários`}
        </Typography>
      </Box>
      <Grid container direction={"row"}>
        {friends.map((friend, idx) => {
          return (
            <Grid
              size={12}
              key={idx}
              sx={{ backgroundColor: isStripedRow(idx, 1) ? palette.darkBlue : palette.mediumBlue }}
            >
              <Stack alignItems={"center"} flexDirection={"row"} sx={{ textAlign: "left" }} gap={1}>
                <img height={64} src={friend.profilePicture} alt="image not found" />
                <Typography variant={"body2"}>{`${friend.name} ${friend.birthDay.toDateString()}`}</Typography>
              </Stack>
            </Grid>
          );
        })}
      </Grid>
    </Paper>
  );
};
