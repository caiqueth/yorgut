import React, { useEffect, useState } from "react";
import { LateralButtonList } from "components/molecules/LateralButtonList";

import { Box } from "@mui/material";
import {
  useJoinCommunity,
  useLeaveCommunity,
  useUserMemberships,
} from "api/user/queries";
import { useNavigate } from "react-router-dom";

interface IProfileMenuProps {
  communityId: number;
}
export const ProfileMenu: React.FC<IProfileMenuProps> = ({ communityId }) => {
  const userId  = 1
  const [isMemberOf, setIsMemberOf] = useState(false);
  const { data: memberships, isFetching } = useUserMemberships();
  const { mutate: joinCommunity } = useJoinCommunity();
  const { mutate: leaveCommunity } = useLeaveCommunity();

  const navigate = useNavigate();
  const handleGoToCommunity = () => navigate(`/community/${communityId}`);
  const handleJoinCommunity = () => joinCommunity({ userId, communityId });
  const handleLeaveCommunity = () => leaveCommunity({ userId, communityId });

  useEffect(() => {
    setIsMemberOf(
      memberships?.communities.some((x) => x.id === communityId) ?? false,
    );
  }, [isFetching]);

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
