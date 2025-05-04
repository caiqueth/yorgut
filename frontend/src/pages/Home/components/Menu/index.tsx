import React from "react";
import { LateralButtonList } from "components/molecules/LateralButtonList";
import { useNavigate } from "react-router-dom";

import { Box } from "@mui/material";

interface IHomeMenuProps {
}
export const Menu: React.FC<IHomeMenuProps> = () => {
  const navigate = useNavigate();

  const handleGoToProfile = () => {
    navigate(`/profile`);
  };

  const handleGoToAlbums = () => {
    navigate(`/profile/albums`);
  };

  const handleGoToScraps = () => {
    navigate(`/profile/scraps`);
  };

  const handleGoToTestimonials = () => {
    navigate(`/profile/testimonials`);
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
