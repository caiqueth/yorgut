import React from "react";
import { LateralButtonList } from "components/molecules/LateralButtonList";
import { useNavigate } from "react-router-dom";

import { Box } from "@mui/material";
import { useMe } from "api/me/queries";

interface IHomeMenuProps {
}
export const Menu: React.FC<IHomeMenuProps> = () => {
  const navigate = useNavigate();
  const { data: user } = useMe();

  if (!user) return null;

  const handleGoToProfile = () => {
    navigate(`/profile/${user.id}`);
  };

  const handleGoToAlbums = () => {
    navigate(`/profile/${user.id}/albums`);
  };

  const handleGoToScraps = () => {
    navigate(`/profile/${user.id}/scraps`);
  };

  const handleGoToTestimonials = () => {
    navigate(`/profile/${user.id}/testimonials`);
  };

  const handleGoToMyAccount = () => {
    navigate(`/my-account`);
  };

  const listItems = [
    { text: "perfil", avatar: "🙂", onClick: handleGoToProfile },
    { text: "álbum", avatar: "📷", onClick: handleGoToAlbums },
    { text: "listas", avatar: "📃" },
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
