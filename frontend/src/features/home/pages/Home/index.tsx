import { CommunitiesList } from "components/CommunitiesList";
import { FriendsList } from "components/FriendsList";
import { UserProfilePicture } from "components/UserProfilePicture";

import { Box, Grid2 as Grid, Stack } from "@mui/material";

import { Birthdays } from "features/home/components/Birthdays";
import { HomeMenu } from "components/HomeMenu";
import { Welcome } from "features/home/components/Welcome";
import { PendingInvites } from "features/home/components/PendingInvites";
import { useMe } from "api/me/queries";


export const Home = () => {
  const { data: user } = useMe();

  if (!user) return null;

  return (
    <Box sx={{ width: "100%" }}>
      <Grid container spacing={2}>
        <Grid size={2}>
          <Stack direction={"column"} gap={2}>
            <UserProfilePicture userId={user.id} />
            <HomeMenu />
          </Stack>
        </Grid>
        <Grid size={10}>
          <Grid container spacing={2}>
            <Grid size={7}>
              <Stack direction={"column"} gap={2}>
                <Welcome />
                <PendingInvites />
                <Birthdays />
              </Stack>
            </Grid>
            <Grid size={5}>
              <Stack direction={"column"} gap={2}>
                <FriendsList />
                <CommunitiesList />
              </Stack>
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </Box>
  );
};
