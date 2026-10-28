import React, { useState } from "react";
import { ILateralButtonItem, LateralButtonList, placeholderItem } from "components/LateralButtonList";

import { Box } from "@mui/material";
import { useNavigate } from "react-router-dom";
import {
  useCancelJoinRequest,
  useJoinCommunity,
  useLeaveCommunity,
  useMe,
  useMeAmMember,
} from "api/me/queries";
import { useCommunity, useCommunityJoinRequests } from "api/community/queries";
import { ConfirmDialog } from "components/ConfirmDialog";

interface IProfileMenuProps {
  communityId: number;
}
export const ProfileMenu: React.FC<IProfileMenuProps> = ({ communityId }) => {
  const { data } = useMeAmMember(communityId);
  const { data: me } = useMe();
  const { data: community } = useCommunity(communityId);
  const { mutate: joinCommunity, isPending: isJoining } = useJoinCommunity();
  const { mutate: leaveCommunity, isPending: isLeaving } = useLeaveCommunity();
  const { mutate: cancelJoinRequest, isPending: isCancelling } = useCancelJoinRequest();
  // "deixar comunidade" asks first.
  const [confirmingLeave, setConfirmingLeave] = useState(false);

  const isMemberOf = data?.amIMember;
  const isPending = data?.pending;

  // The owner is always a member and can't leave, so they get no join/leave option.
  const isOwner = !!me && !!community && me.id === community.ownerId;
  const { data: joinRequests } = useCommunityJoinRequests(communityId, isOwner && !!community?.private);

  const navigate = useNavigate();
  const handleGoToCommunity = () => navigate(`/community/${communityId}`);
  const handleGoToForum = () => navigate(`/community/${communityId}/forum`);
  const handleGoToMembers = () => navigate(`/community/${communityId}/members`);
  const handleGoToEdit = () => navigate(`/community/${communityId}/edit`);
  const handleGoToRequests = () => navigate(`/community/${communityId}/requests`);
  const handleGoToInvite = () => navigate(`/community/${communityId}/invite`);
  const handleGoToPolls = () => navigate(`/community/${communityId}/polls`);
  const handleJoinCommunity = () => joinCommunity({ communityId });
  const handleLeaveCommunity = () => leaveCommunity({ communityId });

  const membershipItem: ILateralButtonItem = isMemberOf
    ? { text: "deixar comunidade", avatar: "➖", onClick: () => setConfirmingLeave(true), loading: isLeaving }
    : isPending
      ? {
        text: "cancelar pedido",
        avatar: "⏳",
        onClick: () => cancelJoinRequest({ communityId }),
        loading: isCancelling,
      }
      : {
        text: community?.private ? "pedir para participar" : "participar",
        avatar: "➕",
        onClick: () => handleJoinCommunity(),
        loading: isJoining,
      };
  // Until we know whether the user owns the community or is a member, the
  // join/leave line holds its place.
  const membershipItems = !me || !community || (!isOwner && !data)
    ? [placeholderItem]
    : isOwner ? [] : [membershipItem];

  const itemsList: ILateralButtonItem[] = [
    { text: "perfil", avatar: "🏘️", onClick: () => handleGoToCommunity() },
    ...membershipItems,
    // Only members may invite their friends.
    ...(isMemberOf ? [{ text: "convidar amigos", avatar: "📧", onClick: () => handleGoToInvite() }] : []),
    { text: "fórum", avatar: "📫", onClick: () => handleGoToForum() },
    { text: "enquetes", avatar: "📔", onClick: () => handleGoToPolls() },
    { text: "membros", avatar: "👥", onClick: () => handleGoToMembers() },
  ];

  if (isOwner && community?.private) {
    itemsList.unshift({
      // No count until it loads, rather than a wrong 0.
      text: `pedidos para entrar${joinRequests ? ` (${joinRequests.requests.length})` : ""}`,
      avatar: "📨",
      onClick: () => handleGoToRequests(),
    });
  }
  if (isOwner) {
    itemsList.unshift({ text: "editar comunidade", avatar: "✏️", onClick: () => handleGoToEdit() });
  }

  return (
    <Box>
      <LateralButtonList items={itemsList} />
      <ConfirmDialog
        open={confirmingLeave}
        title={"deixar comunidade"}
        text={`Tem certeza que quer sair de ${community?.name ?? "esta comunidade"}?`}
        confirmLabel={"sair"}
        onConfirm={handleLeaveCommunity}
        onClose={() => setConfirmingLeave(false)}
      />
    </Box>
  );
};
