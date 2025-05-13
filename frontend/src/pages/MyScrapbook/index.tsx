import { useMe } from "api/me/queries";
import React from "react";
import { Scrapbook } from "components/organisms/Scrapbook";
import { Box, Grid2 as Grid, Stack } from "@mui/material";
import { UserProfilePicture } from "components/molecules/ProfilePicture";
import { Menu } from "pages/Home/components/Menu";

interface IProps { }
export const MyScrapbook: React.FC<IProps> = () => {
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
        <Grid size={10}>
          <Scrapbook userId={user.id} />
        </Grid>
      </Grid>
    </Box>
  );
};
