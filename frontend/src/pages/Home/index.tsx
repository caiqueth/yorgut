import { CommunitiesList } from "./components/CommunitiesList";
import { FriendsList } from "./components/FriendsList";
import { UserProfilePicture } from "components/organisms/UserProfilePicture";

import { Box, Grid2 as Grid, Stack } from "@mui/material";

import { Birthdays } from "./components/Birthdays";
import { Menu } from "./components/Menu";
import { Welcome } from "./components/Welcome";
import { PendingInvites } from "./components/PendingInvites";
import { useLoggedUser } from "api/user/queries";

export const Home = () => {
  const { userId } = useLoggedUser();

  return (
    <Box sx={{ width: "100%" }}>
      <Grid container spacing={2}>
        <Grid size={2}>
          <Stack direction={"column"} gap={2}>
            <UserProfilePicture userId={userId} />
            <Menu userId={userId} />
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
