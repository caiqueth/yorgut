import { useIdParam } from "hooks/useIdParam";
import { Grid2 as Grid, Stack, useTheme } from "@mui/material";
import React from "react";
import { ProfileHeader } from "features/profile/components/ProfileHeader";
import { FriendsList } from "components/FriendsList";
import { CommunitiesList } from "components/CommunitiesList";
import { ProfileDetails } from "features/profile/components/ProfileDetails";
import { TestimonialsPanel } from "features/profile/components/TestimonialsPanel";

export const MainPage: React.FC = () => {
  const userId = useIdParam("userId");
  // Some themes draw the header and the fields as one box.
  const joined = useTheme().yorgut.joinedProfileHeader;

  return <Grid container spacing={2}>
    <Grid size={{ xs: 12, sm: 7 }}>
      <Stack direction={"column"} gap={2}>
        <Stack direction={"column"} gap={joined ? 0 : 2}>
          <ProfileHeader userId={userId} joinedBelow={joined} />
          <ProfileDetails userId={userId} joinedAbove={joined} />
        </Stack>
        <TestimonialsPanel userId={userId} />
      </Stack>
    </Grid>
    <Grid size={{ xs: 12, sm: 5 }}>
      <Stack direction={"column"} gap={2}>
        <FriendsList userId={userId} />
        <CommunitiesList userId={userId} />
      </Stack>
    </Grid>
  </Grid>
}
