import React from "react";

import { Box, Grid2 as Grid, Stack } from "@mui/material";

import { ProfileMenu } from "features/community/components/ProfileMenu";
import { Outlet, useParams } from "react-router-dom";
import { CommunityProfilePicture } from "features/community/components/CommunityProfilePicture";

export const Community: React.FC = () => {
  const params = useParams();
  const communityId = parseInt(params.communityId ?? "");

  return (
    <Box sx={{ width: "100%" }}>
      <Grid container spacing={2}>
        <Grid size={2}>
          <Stack direction={"column"} gap={2}>
            <CommunityProfilePicture communityId={communityId} />
            <ProfileMenu communityId={communityId} />
          </Stack>
        </Grid>
        <Grid size={10}>
          <Outlet />
        </Grid>
      </Grid>
    </Box>
  );
};
