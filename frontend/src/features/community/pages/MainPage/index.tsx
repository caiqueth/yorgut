import { useIdParam } from "hooks/useIdParam";
import { Grid2 as Grid, Stack, useTheme } from "@mui/material";
import React from "react";
import { ProfileHeader } from "features/community/components/ProfileHeader";
import { ProfileDetails } from "features/community/components/ProfileDetails";
import { MembersList } from "features/community/components/MembersList";
import { RelatedCommunitiesList } from "features/community/components/RelatedCommunitiesList";
import { HotTopics } from "features/community/components/HotTopics";

export const MainPage: React.FC = () => {
  const communityId = useIdParam("communityId");
  // Some themes draw the header and the fields as one box.
  const joined = useTheme().yorgut.joinedProfileHeader;

  return <Grid container spacing={2}>
    <Grid size={{ xs: 12, sm: 7 }}>
      <Stack direction={"column"} gap={2}>
        <Stack direction={"column"} gap={joined ? 0 : 2}>
          <ProfileHeader communityId={communityId} joinedBelow={joined} />
          <ProfileDetails communityId={communityId} joinedAbove={joined} />
        </Stack>
        <HotTopics communityId={communityId} />
      </Stack>
    </Grid>
    <Grid size={{ xs: 12, sm: 5 }}>
      <Stack direction={"column"} gap={2}>
        <MembersList communityId={communityId} />
        <RelatedCommunitiesList communityId={communityId} />
      </Stack>
    </Grid>
  </Grid>
}