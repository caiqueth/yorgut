import React from "react";

import { Box, Grid2 as Grid, Stack } from "@mui/material";

import { ProfileHeader } from "./ProfileHeader";
import { ProfileMenu } from "./ProfileMenu";
import { Outlet, useParams } from "react-router-dom";
import { MembersList } from "components/organisms/MembersList";
import { RelatedCommunitiesList } from "pages/Community/components/RelatedCommunitiesList";
import { CommunityProfilePicture } from "./components/ProfilePicture";

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
        <Grid size={6}>
          <Stack direction={"column"} gap={2}>
            <ProfileHeader communityId={communityId} />
            <Outlet />
          </Stack>
        </Grid>
        <Grid size={4}>
          <Stack direction={"column"} gap={2}>
            <MembersList communityId={communityId} />
            <RelatedCommunitiesList communityId={communityId} />
          </Stack>
        </Grid>
      </Grid>
    </Box>
  );
};
