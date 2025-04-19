import React from "react";

import { Box, Grid2 as Grid, Stack } from "@mui/material";

import { CommunitiesList } from "../../components/organisms/CommunitiesList";
import { FriendsList } from "../../components/organisms/FriendsList";
import { UserProfilePicture } from "../../components/organisms/UserProfilePicture";

import { ProfileDetails } from "./ProfileDetails";
import { ProfileHeader } from "./ProfileHeader";
import { ProfileMenu } from "./ProfileMenu";
import { useParams } from "react-router-dom";

export const Profile: React.FC = () => {
  const params = useParams();
  const userId = parseInt(params.userId ?? "");

  return (
    <Box sx={{ width: "100%" }}>
      <Grid container spacing={2}>
        <Grid size={2}>
          <Stack direction={"column"} gap={2}>
            <UserProfilePicture userId={userId} />
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
