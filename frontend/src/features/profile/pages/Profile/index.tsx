import React from "react";

import { Box, Grid2 as Grid, Stack } from "@mui/material";

import { UserProfilePicture } from "components/UserProfilePicture";

import { ProfileMenu } from "features/profile/components/ProfileMenu";
import { Outlet, useParams } from "react-router-dom";

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
        <Grid size={10}>
          <Outlet />
        </Grid>
      </Grid>
    </Box>
  );
};
