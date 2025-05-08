import { Box, Stack, Typography } from "@mui/material";
import { Button } from "components/atoms/Button";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import React from "react";
import { palette } from "theme";
import { isStripedRow } from "helpers/isStripedRow";
import { Link } from "react-router-dom";
import { useMeFriendsPending } from "api/me/queries";

interface IPendingInvitesProps { };
export const PendingInvites: React.FC<IPendingInvitesProps> = () => {
  const { data } = useMeFriendsPending();

  if (!data?.friends) return null;

  return (
    <BoxWithTitle
      title={`novos pedidos de amizade (${data.friends.length})`}
    >
      {data.friends.map((pf, idx) => (
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
                  {pf.firstName}
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
