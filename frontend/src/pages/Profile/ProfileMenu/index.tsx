import React from "react";
import { LateralButtonList } from "components/molecules/LateralButtonList";

import { Box } from "@mui/material";

interface IProfileMenuProps {
  userId: number;
}
export const ProfileMenu: React.FC<IProfileMenuProps> = ({ userId }) => {
  const itemsList = [
    { text: "recados", avatar: "📝", onClick: () => null },
    { text: "álbum", avatar: "📷", onClick: () => null },
    { text: "vídeos", avatar: "📹", onClick: () => null },
    { text: "enviar mensagem", avatar: "📧", onClick: () => null },
    { text: "escrever depoimento", avatar: "✒", onClick: () => null },
    { text: "provocar", avatar: "😋", onClick: () => null },
    { text: "ignorar", avatar: "🤐", onClick: () => null },
    { text: "reportar", avatar: "🚫", onClick: () => null },
  ];

  return (
    <Box>
      <LateralButtonList items={itemsList} />
    </Box>
  );
};
