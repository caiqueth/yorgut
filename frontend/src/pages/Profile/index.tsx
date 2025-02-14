import React from "react";

import { Box, Grid2 as Grid, Stack } from "@mui/material";

import { CommunitiesList } from "../../components/CommunitiesList";
import { FriendsList } from "../../components/FriendsList";
import { ProfilePicture } from "../../components/ProfilePicture";

import { ProfileDetails } from "./ProfileDetails";
import { ProfileHeader } from "./ProfileHeader";
import { ProfileMenu } from "./ProfileMenu";

export const Profile: React.FC = () => {
  const userId = 1;
  return (
    <Box sx={{ width: "100%" }}>
      <Grid container spacing={2}>
        <Grid size={2}>
          <Stack direction={"column"} gap={2}>
            <ProfilePicture userId={userId} />
            <ProfileMenu userId={userId} />
          </Stack>
        </Grid>
        <Grid size={6}>
          <Stack direction={"column"} gap={2}>
            <ProfileHeader userId={userId} />
            <ProfileDetails userId={userId} />
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
