import React, { useEffect, useState } from "react";
import { LateralButtonList } from "components/molecules/LateralButtonList";

import { Box } from "@mui/material";
import {
  useJoinCommunity,
  useLeaveCommunity,
  useLoggedUser,
  useUserMemberships,
} from "api/user/queries";

interface IProfileMenuProps {
  communityId: number;
}
export const ProfileMenu: React.FC<IProfileMenuProps> = ({ communityId }) => {
  const { userId } = useLoggedUser();
  const [isMemberOf, setIsMemberOf] = useState(false);
  const { data: memberships, isPending } = useUserMemberships(userId);
  const { mutate: joinCommunity } = useJoinCommunity();
  const { mutate: leaveCommunity } = useLeaveCommunity();

  const handleJoinCommunity = () => joinCommunity({ userId, communityId });
  const handleLeaveCommunity = () => leaveCommunity({ userId, communityId });

  useEffect(() => {
    setIsMemberOf(
      memberships?.communities.some((x) => x.id === communityId) ?? false,
    );
  }, [isPending]);

  const itemsList = [
    {
      text: isMemberOf ? "deixar comunidade" : "juntar-se à comunidade",
      avatar: "➖",
      onClick: () =>
        isMemberOf ? handleLeaveCommunity() : handleJoinCommunity(),
    },
    { text: "convidar amigos", avatar: "➕", onClick: () => null },
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
