import React from "react";
import { LateralButtonList } from "components/molecules/LateralButtonList";

import { Box } from "@mui/material";
import { useNavigate } from "react-router-dom";
import { useJoinCommunity, useLeaveCommunity, useMe, useMeAmMember } from "api/me/queries";
import { useCommunity } from "api/community/queries";

interface IProfileMenuProps {
  communityId: number;
}
export const ProfileMenu: React.FC<IProfileMenuProps> = ({ communityId }) => {
  const { data } = useMeAmMember(communityId);
  const { data: me } = useMe();
  const { data: community } = useCommunity(communityId);
  const { mutate: joinCommunity } = useJoinCommunity();
  const { mutate: leaveCommunity } = useLeaveCommunity();

  const isMemberOf = data?.amIMember;

  const navigate = useNavigate();
  const handleGoToCommunity = () => navigate(`/community/${communityId}`);
  const handleGoToForum = () => navigate(`/community/${communityId}/forum`);
  const handleGoToMembers = () => navigate(`/community/${communityId}/members`);
  const handleGoToEdit = () => navigate(`/community/${communityId}/edit`);
  const handleJoinCommunity = () => joinCommunity({ communityId });
  const handleLeaveCommunity = () => leaveCommunity({ communityId });

  // The owner is always a member and can't leave, so they get no join/leave option.
  const isOwner = !!me && !!community && me.id === community.ownerId;
  const membershipItems = isOwner ? [] : [
    {
      text: isMemberOf ? "deixar comunidade" : "participar",
      avatar: isMemberOf ? "➖" : "➕",
      onClick: () =>
        isMemberOf ? handleLeaveCommunity() : handleJoinCommunity(),
    },
  ];

  const itemsList = [
    { text: "perfil", avatar: "🏘️", onClick: () => handleGoToCommunity() },
    ...membershipItems,
    { text: "convidar amigos", avatar: "📧" },
    { text: "denunciar abuso", avatar: "⚠️" },
    { text: "fórum", avatar: "📫", onClick: () => handleGoToForum() },
    { text: "enquetes", avatar: "📔" },
    { text: "membros", avatar: "🧑‍🤝‍🧑", onClick: () => handleGoToMembers() },
  ];

  if (isOwner) {
    itemsList.unshift({ text: "editar comunidade", avatar: "✏️", onClick: () => handleGoToEdit() });
  }

  return (
    <Box>
      <LateralButtonList items={itemsList} />
    </Box>
  );
};
