import React from "react";
import { useUserProfile } from "api/user/queries";

import { Box } from "@mui/material";

interface IProfilePictureProps {
  userId: number;
}
export const ProfilePicture: React.FC<IProfilePictureProps> = ({ userId }) => {
  const { data: user } = useUserProfile(userId);
  return (
    <Box>
      <Box
        sx={{
          boxShadow: "inset 1px 1px 20px 14px rgba(255,255,255,1)",
          background: `url(/${user?.profilePicture}) no-repeat center top`,
          backgroundSize: "cover",
          display: "block",
          width: "100%",
          height: "15em",
        }}
      ></Box>
    </Box>
  );
};
