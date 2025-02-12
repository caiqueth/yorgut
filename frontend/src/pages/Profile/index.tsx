import React from "react";
import ProfileImage from "assets/users/frizz.jpg";

import { Box, Grid2 as Grid, Stack } from "@mui/material";

import { CommunitiesList } from "./CommunitiesList";
import { FriendsList } from "./FriendsList";
import { ProfileDetails } from "./ProfileDetails";
import { ProfileHeader } from "./ProfileHeader";
import { ProfileMenu } from "./ProfileMenu";
import { ProfilePicture } from "./ProfilePicture";

export const Profile: React.FC = () => {
  const userId = 1;
  return (
    <Box sx={{ width: "100%" }}>
      <Grid container spacing={2}>
        <Grid size={2}>
          <Stack direction={"column"} gap={2}>
            <ProfilePicture imagePath={ProfileImage} />
            <ProfileMenu />
          </Stack>
        </Grid>
        <Grid size={6}>
          <Stack direction={"column"} gap={2}>
            <ProfileHeader />
            <ProfileDetails />
          </Stack>
        </Grid>
        <Grid size={4}>
          <Stack direction={"column"} gap={2}>
            <FriendsList userId={userId} />
            <CommunitiesList userId={userId} />
          </Stack>
        </Grid>
      </Grid>
    </Box>
  );
};
