import React from "react";
import { LateralButtonList } from "components/molecules/LateralButtonList";
import { useNavigate } from "react-router-dom";

import { Box } from "@mui/material";

interface IHomeMenuProps {
  userId: number;
}
export const HomeMenu: React.FC<IHomeMenuProps> = ({ userId }) => {
  const navigate = useNavigate();
  const handleGoToProfile = () => {
    navigate(`/profile/${userId}`);
  };

  const listItems = [
    { text: "perfil", avatar: "🙂", onClick: handleGoToProfile },
    { text: "álbum", avatar: "📷", onClick: () => null },
    { text: "listas", avatar: "📃", onClick: () => null },
    { text: "recados", avatar: "📝", onClick: () => null },
    { text: "depoimentos", avatar: "✒", onClick: () => null },
    { text: "configurações", avatar: "🔨", onClick: () => null },
  ];

  return (
    <Box>
      <LateralButtonList items={listItems} />
    </Box>
  );
};
