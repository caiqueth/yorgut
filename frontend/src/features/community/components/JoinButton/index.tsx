import { Stack, Typography } from "@mui/material";
import { useCommunity } from "api/community/queries";
import { useCancelJoinRequest, useJoinCommunity, useMeAmMember } from "api/me/queries";
import { Button } from "components/Button";
import React from "react";

interface IProps {
  communityId: number;
}
// The way into a community for a non-member: "participar" for a public one,
// "pedir para participar" for a private one, or the pending request and a way
// to cancel it. Renders nothing for members.
export const JoinButton: React.FC<IProps> = ({ communityId }) => {
  const { data: community } = useCommunity(communityId);
  const { data: membership } = useMeAmMember(communityId);
  const { mutate: joinCommunity, isPending: isJoining } = useJoinCommunity();
  const { mutate: cancelJoinRequest, isPending: isCancelling } = useCancelJoinRequest();

  if (!community || !membership || membership.amIMember) return null;

  if (membership.pending) {
    return (
      <Stack direction={"row"} alignItems={"center"} gap={1} flexWrap={"wrap"}>
        <Typography variant={"body2"}>⏳ Pedido enviado. Aguardando o dono aceitar.</Typography>
        <Button onClick={() => cancelJoinRequest({ communityId })} loading={isCancelling}>
          cancelar pedido
        </Button>
      </Stack>
    );
  }

  return (
    <Button onClick={() => joinCommunity({ communityId })} loading={isJoining}>
      {community.private ? "pedir para participar" : "participar"}
    </Button>
  );
};
