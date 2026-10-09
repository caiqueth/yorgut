import { Box, Grid2 as Grid, Stack } from "@mui/material";
import React from "react";
import { ProfileHeader } from "features/profile/components/ProfileHeader";
import { FriendsList } from "features/profile/components/FriendsList";
import { CommunitiesList } from "features/profile/components/CommunitiesList";
import { useParams } from "react-router-dom";
import { ProfileDetails } from "features/profile/components/ProfileDetails";
import { TestimonialsPanel } from "features/profile/components/TestimonialsPanel";

interface IProps { };
export const MainPage: React.FC<IProps> = () => {
  const params = useParams()
  const userId = parseInt(params.userId ?? "");

  return <Grid container spacing={2}>
    <Grid size={7}>
      <Stack direction={"column"} gap={2}>
        <ProfileHeader userId={userId} />
        <ProfileDetails userId={userId} />
        <TestimonialsPanel userId={userId} />
      </Stack>
    </Grid>
    <Grid size={5}>
      <Stack direction={"column"} gap={2}>
        <FriendsList userId={userId} />
        <CommunitiesList userId={userId} />
      </Stack>
    </Grid>
  </Grid>
}
