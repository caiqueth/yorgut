import { Box, Stack, Typography } from "@mui/material";
import { Button } from "components/atoms/Button";
import { useUserPendingFriendships } from "api/user/queries";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import React from "react";
import { palette } from "theme";
import { isStripedRow } from "helpers/isStripedRow";
import { Link } from "react-router-dom";

interface IPendingInvitesProps {
  userId: number;
}
export const PendingInvites: React.FC<IPendingInvitesProps> = ({ userId }) => {
  const { data } = useUserPendingFriendships(userId);

  if (!data) return null;
  if (data.pendingFriendships.length <= 0) return null;

  return (
    <BoxWithTitle
      title={`novos pedidos de amizade (${data.pendingFriendships.length})`}
    >
      {data?.pendingFriendships.map((pf, idx) => (
        <Box
          sx={{
            backgroundColor: isStripedRow(idx, 1)
              ? palette.darkBlue
              : palette.mediumBlue,
            px: 1,
          }}
          key={idx}
        >
          <Stack direction={"row"} alignItems={"center"} gap={1}>
            <Link to={`/profile/${pf.id}`}>
              <Box>
                <img src={`${pf.profilePicture}`} height={64}></img>
              </Box>
            </Link>
            <Box>
              <Stack direction={"column"}>
                <Typography>
                  {pf.firstName} {pf.lastName}
                </Typography>
                <Typography>{pf.firstName} é seu(sua) amigo(a)?</Typography>
                <Box>
                  <Button>Sim</Button>
                  <Button>Não</Button>
                </Box>
              </Stack>
            </Box>
          </Stack>
        </Box>
      ))}
    </BoxWithTitle>
  );
};
