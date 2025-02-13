import { CommunitiesList } from "pages/Profile/CommunitiesList";
import { FriendsList } from "pages/Profile/FriendsList";
import { ProfilePicture } from "pages/Profile/ProfilePicture";

import { Box, Grid2 as Grid, Stack } from "@mui/material";

import { Birthdays } from "./Birthdays";
import { HomeMenu } from "./HomeMenu";
import { Welcome } from "./Welcome";

export const Home = () => {
  const userId = 1;
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
            <Welcome />
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
