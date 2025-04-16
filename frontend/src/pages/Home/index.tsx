import { CommunitiesList } from "components/CommunitiesList";
import { FriendsList } from "components/FriendsList";
import { ProfilePicture } from "components/ProfilePicture";

import { Box, Grid2 as Grid, Stack } from "@mui/material";

import { Birthdays } from "./Birthdays";
import { HomeMenu } from "./HomeMenu";
import { Welcome } from "./Welcome";
import { PendingInvites } from "./PendingInvites";

export const Home = () => {
  const userId = 6;
  return (
    <Box sx={{ width: "100%" }}>
      <Grid container spacing={2}>
        <Grid size={2}>
          <Stack direction={"column"} gap={2}>
            <ProfilePicture userId={userId} />
            <HomeMenu userId={userId} />
          </Stack>
        </Grid>
        <Grid size={6}>
          <Stack direction={"column"} gap={2}>
            <Welcome userId={userId} />
            <PendingInvites userId={userId} />
            <Birthdays userId={userId} />
          </Stack>
        </Grid>
        <Grid size={4}>
          <Stack direction={"column"} gap={2}>
            <FriendsList userId={userId} />
            <CommunitiesList userId={userId} />
          </Stack>
        </Grid>
      </Grid>
    </Box>
  );
};
