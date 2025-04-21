import React from "react";
import { LateralButtonList } from "components/molecules/LateralButtonList";
import { useNavigate } from "react-router-dom";

import { Box } from "@mui/material";

interface IHomeMenuProps {
  userId: number;
}
export const Menu: React.FC<IHomeMenuProps> = ({ userId }) => {
  const navigate = useNavigate();

  const handleGoToProfile = () => {
    navigate(`/profile/${userId}`);
  };

  const handleGoToAlbums = () => {
    navigate(`/profile/${userId}/albums`);
  };

  const handleGoToScraps = () => {
    navigate(`/profile/${userId}/scraps`);
  };

  const handleGoToTestimonials = () => {
    navigate(`/profile/${userId}/testimonials`);
  };

  const handleGoToMyAccount = () => {
    navigate(`/my-account`);
  };

  const listItems = [
    { text: "perfil", avatar: "🙂", onClick: handleGoToProfile },
    { text: "álbum", avatar: "📷", onClick: handleGoToAlbums },
    { text: "listas", avatar: "📃", onClick: () => null },
    { text: "recados", avatar: "📝", onClick: handleGoToScraps },
    {
      text: "depoimentos",
      avatar: "🖊️",
      onClick: handleGoToTestimonials,
    },
    { text: "configurações", avatar: "🔨", onClick: handleGoToMyAccount },
  ];

  return (
    <Box>
      <LateralButtonList items={listItems} />
    </Box>
  );
};
