import React from "react";
import { useUserFriends } from "api/user/queries";
import { isStripedRow } from "helpers/isStripedRow";
import { Link } from "react-router-dom";
import { palette } from "theme";

import { Grid2 as Grid, Stack, Typography } from "@mui/material";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { dayjsFromDate } from "helpers/formatting/date";

interface IBirthdaysProps {
  userId: number;
}
export const Birthdays: React.FC<IBirthdaysProps> = ({ userId }) => {
  const { data: friends } = useUserFriends(userId);
  return (
    <BoxWithTitle title={`próximos aniversários`}>
      <Grid container direction={"row"}>
        {friends?.friends?.map((friend, idx) => {
          return (
            <Grid
              size={12}
              key={idx}
              sx={{
                backgroundColor: isStripedRow(idx, 1)
                  ? palette.darkBlue
                  : palette.mediumBlue,
                p: 1,
              }}
            >
              <Stack
                alignItems={"center"}
                flexDirection={"row"}
                sx={{ textAlign: "left" }}
                gap={1}
              >
                <img
                  height={64}
                  src={friend.profilePicture}
                  alt="image not found"
                />
                <Typography variant={"body2"}>
                  <Link to={`/profile/${friend.id}`}>{friend.firstName}</Link>{" "}
                  {`${dayjsFromDate(friend.birthDate, "D [de] MMMM")}`}
                </Typography>
              </Stack>
            </Grid>
          );
        })}
      </Grid>
    </BoxWithTitle>
  );
};
