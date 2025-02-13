import React from "react";
import { useUserCommunities, useUserStatistics } from "api/user/queries";
import { palette } from "theme";

import { Box, Grid2 as Grid, Paper, Stack, Typography } from "@mui/material";

interface ICommunitiesListProps {
  userId: number;
}
export const CommunitiesList: React.FC<ICommunitiesListProps> = ({ userId }) => {
  const userStatistics = useUserStatistics(userId);
  const communities = useUserCommunities(userId);

  return (
    <Paper>
      <Box sx={{ textAlign: "center", width: "100%", backgroundColor: palette.lightBlue }}>
        <Typography variant={"subtitle2"} fontWeight={"bold"}>
          {`comunidades (${userStatistics.totalCommunities})`}
        </Typography>
      </Box>
      <Grid container spacing={1} sx={{ p: 1 }}>
        {communities.map((community) => (
          <Grid size={4} key={community.id}>
            <Stack alignItems={"center"} flexDirection={"column"} sx={{ textAlign: "center" }}>
              <img height={64} src={community.picture} alt="image not found" />
              <Typography variant={"body2"}>{`${community.name} (${community.members})`}</Typography>
            </Stack>
          </Grid>
        ))}
      </Grid>
    </Paper>
  );
};
