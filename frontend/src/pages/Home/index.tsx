import { CommunitiesList } from "./components/CommunitiesList";
import { FriendsList } from "./components/FriendsList";
import { UserProfilePicture } from "components/molecules/UserProfilePicture";

import { Box, Grid2 as Grid, Stack } from "@mui/material";

import { Birthdays } from "./components/Birthdays";
import { Menu } from "./components/Menu";
import { Welcome } from "./components/Welcome";
import { PendingInvites } from "./components/PendingInvites";
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
            <Menu />
          </Stack>
        </Grid>
        <Grid size={6}>
          <Stack direction={"column"} gap={2}>
            <Welcome />
            <PendingInvites />
            <Birthdays />
          </Stack>
        </Grid>
        <Grid size={4}>
          <Stack direction={"column"} gap={2}>
            <FriendsList />
            <CommunitiesList />
          </Stack>
        </Grid>
      </Grid>
    </Box>
  );
};
