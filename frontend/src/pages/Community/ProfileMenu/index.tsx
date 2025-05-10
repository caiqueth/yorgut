import React from "react";
import { LateralButtonList } from "components/molecules/LateralButtonList";

import { Box } from "@mui/material";
import { useNavigate } from "react-router-dom";
import { useJoinCommunity, useLeaveCommunity, useMeAmMember } from "api/me/queries";

interface IProfileMenuProps {
  communityId: number;
}
export const ProfileMenu: React.FC<IProfileMenuProps> = ({ communityId }) => {
  const { data } = useMeAmMember(communityId);
  const { mutate: joinCommunity } = useJoinCommunity();
  const { mutate: leaveCommunity } = useLeaveCommunity();

  const isMemberOf = data?.amIMember;

  const navigate = useNavigate();
  const handleGoToCommunity = () => navigate(`/community/${communityId}`);
  const handleJoinCommunity = () => joinCommunity({ communityId });
  const handleLeaveCommunity = () => leaveCommunity({ communityId });

  const itemsList = [
    { text: "perfil", avatar: "🏘️", onClick: () => handleGoToCommunity() },
    {
      text: isMemberOf ? "deixar comunidade" : "participar",
      avatar: isMemberOf ? "➖" : "➕",
      onClick: () =>
        isMemberOf ? handleLeaveCommunity() : handleJoinCommunity(),
    },
    { text: "convidar amigos", avatar: "📧", onClick: () => null },
    { text: "denunciar abuso", avatar: "⚠️", onClick: () => null },
    { text: "fórum", avatar: "📫", onClick: () => null },
    { text: "enquetes", avatar: "📔", onClick: () => null },
    { text: "membros", avatar: "🧑‍🤝‍🧑", onClick: () => null },
  ];

  return (
    <Box>
      <LateralButtonList items={itemsList} />
    </Box>
  );
};
