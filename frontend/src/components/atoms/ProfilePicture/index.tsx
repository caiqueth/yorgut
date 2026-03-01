import React from "react";

import { Box, Paper, Skeleton } from "@mui/material";

interface IProfilePictureProps {
  imageSrc?: string;
  isLoading?: boolean;
}
export const ProfilePicture: React.FC<IProfilePictureProps> = ({
  imageSrc,
  isLoading,
}) => {
  return (
    isLoading
      ? (
        <Box
          sx={{
            boxShadow: "inset 1px 1px 20px 12px rgba(255,255,255,1)",
            backgroundSize: "cover",
            display: "block",
            width: "100%",
            height: "15em",
            border: ""
          }}
        >
          <Skeleton height={"100%"} />
        </Box>
      ) : (
        <Box
          sx={{
            boxShadow: "inset 1px 1px 20px 12px rgba(255,255,255,1)",
            background: imageSrc ? `url(${imageSrc}) no-repeat center top` : "",
            backgroundSize: "cover",
            display: "block",
            width: "100%",
            height: "15em",
          }}
        />
      )
  );
};
