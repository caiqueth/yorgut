import React from "react";
import { Stack, Typography } from "@mui/material";
import {
  useAcceptCommunityInvitation,
  useDeclineCommunityInvitation,
  useMyCommunityInvitations,
} from "api/me/queries";
import { BoxWithTitle } from "components/BoxWithTitle";
import { Button } from "components/Button";
import { PicAndTextHorizontal } from "components/PicAndTextHorizontal";
import { StripedGrid } from "components/StripedGrid";

// The invitations friends sent the logged user to join their communities.
// Hidden when there are none.
export const CommunityInvitations: React.FC = () => {
  const { data, isFetching } = useMyCommunityInvitations();
  const { mutate: accept, isPending: isAccepting, variables: accepting } = useAcceptCommunityInvitation();
  const { mutate: decline, isPending: isDeclining, variables: declining } = useDeclineCommunityInvitation();

  if (!data || data.invitations.length === 0) return null;

  const items = data.invitations.map((inv) => {
    const variables = { invitationId: inv.id, communityId: inv.communityId };
    const isAnswering =
      (isAccepting && accepting?.invitationId === inv.id) || (isDeclining && declining?.invitationId === inv.id);
    return (
      <Stack key={inv.id} direction={"row"} justifyContent={"space-between"} alignItems={"center"} gap={1} flexWrap={"wrap"}>
        <PicAndTextHorizontal
          title={`${inv.communityPrivate ? "🔒 " : ""}${inv.communityName}`}
          text={`${inv.inviterFirstName} convidou você para participar`}
          url={`/community/${inv.communityId}`}
          imageSrc={inv.communityPicture}
          imageSize={56}
        />
        <Stack direction={"row"}>
          <Button
            onClick={() => accept(variables)}
            loading={isAccepting && accepting?.invitationId === inv.id}
            disabled={isAnswering}
          >
            participar
          </Button>
          <Button
            onClick={() => decline(variables)}
            loading={isDeclining && declining?.invitationId === inv.id}
            disabled={isAnswering}
          >
            recusar
          </Button>
        </Stack>
      </Stack>
    );
  });

  return (
    <BoxWithTitle title={`convites para comunidades (${data.invitations.length})`} showProgress={isFetching}>
      <StripedGrid items={items} elementsPerRow={1} />
      {data.invitations.some((inv) => inv.communityPrivate) && (
        <Typography variant={"body2"} sx={{ p: 1 }}>
          Numa comunidade privada, aceitar o convite de alguém que não é o dono envia um pedido para o dono aprovar.
        </Typography>
      )}
    </BoxWithTitle>
  );
};
