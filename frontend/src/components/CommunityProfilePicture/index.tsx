import React from "react";
import { useCommunityProfile, useUserProfile } from "api/user/queries";

import { Box } from "@mui/material";

interface ICommunityProfilePictureProps {
  communityId: number;
}
export const CommunityProfilePicture: React.FC<
  ICommunityProfilePictureProps
> = ({ communityId }) => {
  const { data: community } = useCommunityProfile(communityId);

  return (
    <Box>
      <Box
        sx={{
          boxShadow: "inset 1px 1px 20px 14px rgba(255,255,255,1)",
          background: `url(/${community?.profilePicture}) no-repeat center top`,
          backgroundSize: "cover",
          display: "block",
          width: "100%",
          height: "15em",
        }}
      ></Box>
    </Box>
  );
};
