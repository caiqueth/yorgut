import React from "react";

import { Box, Grid2 as Grid, Stack } from "@mui/material";

import { ProfilePicture } from "../../components/UserProfilePicture";

import { ProfileDetails } from "./ProfileDetails";
import { ProfileHeader } from "./ProfileHeader";
import { ProfileMenu } from "./ProfileMenu";
import { useParams } from "react-router-dom";
import { MembersList } from "components/MembersList";
import { RelatedCommunitiesList } from "components/RelatedCommunitiesList";
import { CommunityProfilePicture } from "components/CommunityProfilePicture";

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
            <ProfileDetails communityId={communityId} />
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
