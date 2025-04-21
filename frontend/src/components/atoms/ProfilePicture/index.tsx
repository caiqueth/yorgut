import React from "react";

import { Box } from "@mui/material";

interface IProfilePictureProps {
  imageSrc: string;
}
export const ProfilePicture: React.FC<IProfilePictureProps> = ({
  imageSrc,
}) => {
  return (
    <Box
      sx={{
        boxShadow: "inset 1px 1px 20px 14px rgba(255,255,255,1)",
        background: `url(${imageSrc}) no-repeat center top`,
        backgroundSize: "cover",
        display: "block",
        width: "100%",
        height: "15em",
      }}
    ></Box>
  );
};
