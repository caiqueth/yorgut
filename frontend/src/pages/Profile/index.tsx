import React from "react";

import { Grid2 as Grid } from "@mui/material";

import { CommunitiesList } from "./CommunitiesList";
import { FriendsList } from "./FriendsList";
import { ProfileDetails } from "./ProfileDetails";
import { ProfileHeader } from "./ProfileHeader";
import { ProfileMenu } from "./ProfileMenu";
import { ProfilePicture } from "./ProfilePicture";

export const Profile: React.FC = () => {
  const userId = 1;
  return (
    <Grid container>
      <Grid sx={{ xs: 2 }}>
        <ProfilePicture imagePath="/test" />
        <ProfileMenu />
      </Grid>
      <Grid sx={{ xs: 6 }}>
        <ProfileHeader />
        <ProfileDetails />
      </Grid>
      <Grid sx={{ xs: 4 }}>
        <FriendsList userId={userId} />
        <CommunitiesList userId={userId} />
      </Grid>
    </Grid>
  );
};
