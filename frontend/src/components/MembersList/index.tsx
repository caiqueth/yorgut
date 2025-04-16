import React from "react";
import { useCommunity, useCommunityMembers } from "api/user/queries";
import { Button } from "components/Button";
import { isStripedRow } from "helpers/isStripedRow";
import { Link, useNavigate } from "react-router-dom";
import { palette } from "theme";

import { Box, Grid2 as Grid, Stack, Typography } from "@mui/material";
import { BoxWithTitle } from "components/BoxWithTitle";

interface IMembersListProps {
  communityId: number;
}
export const MembersList: React.FC<IMembersListProps> = ({ communityId }) => {
  const { data: members } = useCommunityMembers(communityId);
  const navigate = useNavigate();

  const handleGoToMembers = () => {
    navigate("members");
  };

  return (
    <BoxWithTitle title={`membros (${members?.count})`}>
      <Stack direction={"column"} alignItems={"center"}>
        <Grid container direction={"row"}>
          {members?.members.map((member, idx) => (
            <Grid
              size={3}
              key={idx}
              sx={{
                p: 1,
                backgroundColor: isStripedRow(idx, 4)
                  ? palette.darkBlue
                  : palette.mediumBlue,
              }}
            >
              <Stack
                alignItems={"center"}
                flexDirection={"column"}
                sx={{ textAlign: "center" }}
              >
                <Link to={`/profile/${member.id}`}>
                  <img
                    height={64}
                    src={`/${member.profilePicture}`}
                    alt="image not found"
                  />
                  <Typography
                    variant={"body2"}
                  >{`${member.firstName} (${member.friendCount})`}</Typography>
                </Link>
              </Stack>
            </Grid>
          ))}
        </Grid>
        <Box alignSelf={"center"}>
          <Button onClick={handleGoToMembers}>ver membros</Button>
        </Box>
      </Stack>
    </BoxWithTitle>
  );
};
