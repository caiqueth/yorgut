import React from "react";
import { useUserMemberships, useUserStatistics } from "api/user/queries";
import { isStripedRow } from "helpers/isStripedRow";
import { Link } from "react-router-dom";
import { palette } from "theme";

import { Box, Grid2 as Grid, Paper, Stack, Typography } from "@mui/material";

interface ICommunitiesListProps {
  userId: number;
}
export const CommunitiesList: React.FC<ICommunitiesListProps> = ({ userId }) => {
  const userStatistics = useUserStatistics(userId);
  const { data: memberships } = useUserMemberships(userId);

  return (
    <Paper>
      <Box sx={{ textAlign: "center", width: "100%", backgroundColor: palette.lightBlue }}>
        <Typography variant={"subtitle2"} fontWeight={"bold"}>
          {`comunidades (${userStatistics.totalCommunities})`}
        </Typography>
      </Box>
      <Grid container direction={"row"}>
        {memberships?.communities.map((membership, idx) => (
          <Grid
            size={4}
            key={idx}
            sx={{ p: 1, backgroundColor: isStripedRow(idx, 3) ? palette.darkBlue : palette.mediumBlue }}
          >
            <Stack alignItems={"center"} flexDirection={"column"} sx={{ textAlign: "center" }}>
              <Link to={`/communities/${membership.id}`}>
                <img height={64} src={`/${membership.profilePicture}`} alt="image not found" />
                <Typography variant={"body2"}>{`${membership.name} (${membership.members ?? 0})`}</Typography>
              </Link>
            </Stack>
          </Grid>
        ))}
      </Grid>
    </Paper>
  );
};
