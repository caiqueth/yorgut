import React from "react";

import { Box } from "@mui/material";

interface IProfilePictureProps {
  imagePath: string;
}
export const ProfilePicture: React.FC<IProfilePictureProps> = ({ imagePath }) => {
  return (
    <Box>
      <Box
        sx={{
          boxShadow: "inset 1px 1px 28px 12px rgba(255,255,255,1)",
          background: `url(${imagePath}) no-repeat center top`,
          backgroundSize: "cover",
          display: "block",
          width: "100%",
          height: "15em",
        }}
      ></Box>
    </Box>
  );
};
