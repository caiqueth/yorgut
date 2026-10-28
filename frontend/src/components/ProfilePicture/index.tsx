import React from "react";

import { Box, Skeleton } from "@mui/material";
import { Link } from "react-router-dom";

interface IProfilePictureProps {
  imageSrc?: string | null;
  isLoading?: boolean;
  // Where the picture links to, if anywhere.
  to?: string;
}

// Shorter on phones, where it sits next to the menu instead of above it.
const placeholderHeight = { xs: "11em", md: "15em" };

// The big picture at the top of a profile or community, in its own aspect
// ratio and faded at the edges.
export const ProfilePicture: React.FC<IProfilePictureProps> = ({ imageSrc, isLoading, to }) => {
  const picture = isLoading ? (
    <Skeleton variant={"rectangular"} sx={{ height: placeholderHeight }} />
  ) : (
    <Box
      sx={{
        position: "relative",
        height: imageSrc ? "auto" : placeholderHeight,
        // Over the image, which would hide an inset shadow on the box itself.
        "&::after": {
          content: '""',
          position: "absolute",
          inset: 0,
          boxShadow: "inset 1px 1px 20px 12px rgba(255,255,255,1)",
        },
      }}
    >
      {imageSrc && (
        <Box component={"img"} src={imageSrc} alt={""} sx={{ display: "block", width: "100%", height: "auto" }} />
      )}
    </Box>
  );

  return to ? (
    <Link to={to} style={{ display: "block" }}>
      {picture}
    </Link>
  ) : (
    picture
  );
};
