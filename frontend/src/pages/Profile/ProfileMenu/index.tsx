import React from "react";
import { LateralButtonList } from "components/molecules/LateralButtonList";

import { Box } from "@mui/material";
import { useNavigate } from "react-router-dom";

interface IProfileMenuProps {
  userId: number;
}
export const ProfileMenu: React.FC<IProfileMenuProps> = ({ userId }) => {
  const navigate = useNavigate();

  const handleGoToRoute = (route: string) => {
    navigate(route);
  };

  const itemsList = [
    {
      text: "recados", avatar: "📝", onClick: () => handleGoToRoute(`/profile/${userId}/scraps`)
    },
    { text: "álbum", avatar: "📷", onClick: () => handleGoToRoute("albums") },
    { text: "vídeos", avatar: "📹", onClick: () => handleGoToRoute("videos") },
    {
      text: "escrever depoimento",
      avatar: "🖊️",
      onClick: () => handleGoToRoute("testimonials"),
    },
    { text: "reportar", avatar: "🚫", onClick: () => null },
  ];

  return (
    <Box>
      <LateralButtonList items={itemsList} />
    </Box>
  );
};
