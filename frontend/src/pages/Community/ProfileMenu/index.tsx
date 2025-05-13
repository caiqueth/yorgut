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
  const handleGoToForum = () => navigate(`/community/${communityId}/forum`);
  const handleGoToMembers = () => navigate(`/community/${communityId}/members`);
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
    { text: "convidar amigos", avatar: "📧" },
    { text: "denunciar abuso", avatar: "⚠️" },
    { text: "fórum", avatar: "📫", onClick: () => handleGoToForum() },
    { text: "enquetes", avatar: "📔" },
    { text: "membros", avatar: "🧑‍🤝‍🧑", onClick: () => handleGoToMembers() },
  ];

  return (
    <Box>
      <LateralButtonList items={itemsList} />
    </Box>
  );
};
