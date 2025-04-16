import React from "react";
import { LateralButtonList } from "components/LateralButtonList";

import { Box } from "@mui/material";

interface IProfileMenuProps {
  communityId: number;
}
export const ProfileMenu: React.FC<IProfileMenuProps> = ({ communityId }) => {
  const itemsList = [
    { text: "deixar comunidade", avatar: "➖", onClick: () => null },
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
