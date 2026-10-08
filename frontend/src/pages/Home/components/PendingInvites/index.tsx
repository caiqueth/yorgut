import { Box, Stack, Typography } from "@mui/material";
import { Button } from "components/atoms/Button";
import { uploadUrl } from "helpers/uploadUrl";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import React from "react";
import { palette } from "theme";
import { isStripedRow } from "helpers/isStripedRow";
import { Link } from "react-router-dom";
import { useMeFriendsPending } from "api/me/queries";
import { useReplyFriendship } from "api/friendship/queries";

interface IPendingInvitesProps { };
export const PendingInvites: React.FC<IPendingInvitesProps> = () => {
  const { data } = useMeFriendsPending();
  const { mutate: reply } = useReplyFriendship();

  const handleAcceptFriendship = (friendshipId: number) => {
    reply({
      friendshipId,
      accept: true,
    })
  };

  const handleRejectFriendship = (friendshipId: number) => {
    reply({
      friendshipId,
      accept: false,
    })
  };

  if (!data || data?.count === 0) return null;

  return (
    <BoxWithTitle
      title={`novos pedidos de amizade (${data.count})`}
    >
      {data.pendingFriendships.map((pf, idx) => (
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
            <Link to={`/profile/${pf.userRequestedId}`}>
              <Box>
                <img src={uploadUrl(pf.profilePicture)} height={64}></img>
              </Box>
            </Link>
            <Box>
              <Stack direction={"column"}>
                <Typography>
                  {pf.firstName}
                </Typography>
                <Typography>{pf.firstName} é seu(sua) amigo(a)?</Typography>
                <Box>
                  <Button onClick={() => handleAcceptFriendship(pf.id)}>Sim</Button>
                  <Button onClick={() => handleRejectFriendship(pf.id)}>Não</Button>
                </Box>
              </Stack>
            </Box>
          </Stack>
        </Box >
      ))}
    </BoxWithTitle >
  );
};
