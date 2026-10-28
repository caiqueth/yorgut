import { Box, Stack, Typography } from "@mui/material";
import { Button } from "components/Button";
import { BoxWithTitle } from "components/BoxWithTitle";
import React from "react";
import { isStripedRow } from "helpers/isStripedRow";
import { Link } from "react-router-dom";
import { useMeFriendsPending } from "api/me/queries";
import { useReplyFriendship } from "api/friendship/queries";

interface IPendingInvitesProps { };
export const PendingInvites: React.FC<IPendingInvitesProps> = () => {
  const { data, isFetching } = useMeFriendsPending();
  const { mutate: reply, isPending, variables: replying } = useReplyFriendship();

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
      showProgress={isFetching}
    >
      {data.pendingFriendships.map((pf, idx) => (
        <Box
          sx={{
            backgroundColor: isStripedRow(idx, 1)
              ? "yorgut.rows.striped"
              : "yorgut.rows.plain",
            px: 1,
          }}
          key={idx}
        >
          <Stack direction={"row"} alignItems={"center"} gap={1}>
            <Link to={`/profile/${pf.userRequestedId}`}>
              <Box>
                <img src={pf.profilePicture ?? undefined} height={64}></img>
              </Box>
            </Link>
            <Box>
              <Stack direction={"column"}>
                <Typography>
                  <Link to={`/profile/${pf.userRequestedId}`}>{pf.firstName}</Link>
                </Typography>
                <Typography>{pf.firstName} é seu(sua) amigo(a)?</Typography>
                <Box>
                  <Button
                    onClick={() => handleAcceptFriendship(pf.id)}
                    loading={isPending && replying?.friendshipId === pf.id && replying.accept}
                    disabled={isPending && replying?.friendshipId === pf.id}
                  >
                    Sim
                  </Button>
                  <Button
                    onClick={() => handleRejectFriendship(pf.id)}
                    loading={isPending && replying?.friendshipId === pf.id && !replying.accept}
                    disabled={isPending && replying?.friendshipId === pf.id}
                  >
                    Não
                  </Button>
                </Box>
              </Stack>
            </Box>
          </Stack>
        </Box >
      ))}
    </BoxWithTitle >
  );
};
