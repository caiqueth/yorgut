import React from "react";
import {
  useCommunity,
  useRelatedCommunities,
  useUserMemberships,
  useUserStatistics,
} from "api/user/queries";
import { isStripedRow } from "helpers/isStripedRow";
import { Link } from "react-router-dom";
import { palette } from "theme";

import { Box, Grid2 as Grid, Paper, Stack, Typography } from "@mui/material";
import { BoxWithTitle } from "components/BoxWithTitle";

interface IRelatedCommunitiesListProps {
  communityId: number;
}
export const RelatedCommunitiesList: React.FC<IRelatedCommunitiesListProps> = ({
  communityId,
}) => {
  const { data: communities } = useRelatedCommunities(communityId);

  return (
    <BoxWithTitle
      title={`comunidades relacionadas (${communities?.memberCount})`}
    >
      <Grid container direction={"row"}>
        {communities?.communities.map((community, idx) => (
          <Grid
            size={4}
            key={idx}
            sx={{
              p: 1,
              backgroundColor: isStripedRow(idx, 3)
                ? palette.darkBlue
                : palette.mediumBlue,
            }}
          >
            <Stack
              alignItems={"center"}
              flexDirection={"column"}
              sx={{ textAlign: "center" }}
            >
              <Link to={`/community/${community.id}`}>
                <img
                  height={64}
                  src={`/${community.profilePicture}`}
                  alt="image not found"
                />
                <Typography
                  variant={"body2"}
                >{`${community.name} (${community.membershipCount ?? 0})`}</Typography>
              </Link>
            </Stack>
          </Grid>
        ))}
      </Grid>
    </BoxWithTitle>
  );
};
