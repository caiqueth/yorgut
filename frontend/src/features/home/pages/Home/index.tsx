import React from "react";
import { Grid2 as Grid, Stack } from "@mui/material";
import { useMe } from "api/me/queries";
import { CommunitiesList } from "components/CommunitiesList";
import { FriendsList } from "components/FriendsList";
import { HomeMenu } from "components/HomeMenu";
import { PageWithSidebar } from "components/PageWithSidebar";
import { UserProfilePicture } from "components/UserProfilePicture";
import { Birthdays } from "features/home/components/Birthdays";
import { CommunityInvitations } from "features/home/components/CommunityInvitations";
import { PendingInvites } from "features/home/components/PendingInvites";
import { PendingTestimonials } from "features/home/components/PendingTestimonials";
import { Welcome } from "features/home/components/Welcome";

export const Home: React.FC = () => {
  // Every box loads on its own; only the picture needs the user's id.
  const { data: user } = useMe();

  return (
    <PageWithSidebar picture={<UserProfilePicture userId={user?.id} />} menu={<HomeMenu />}>
      <Grid container spacing={2}>
        <Grid size={{ xs: 12, sm: 7 }}>
          <Stack direction={"column"} gap={2}>
            <Welcome />
            <PendingInvites />
            <CommunityInvitations />
            <PendingTestimonials />
            <Birthdays />
          </Stack>
        </Grid>
        <Grid size={{ xs: 12, sm: 5 }}>
          <Stack direction={"column"} gap={2}>
            <FriendsList />
            <CommunitiesList />
          </Stack>
        </Grid>
      </Grid>
    </PageWithSidebar>
  );
};
