import React from "react";
import { Box, Skeleton, Stack, Typography } from "@mui/material";
import { Link } from "react-router-dom";

const IMAGE_HEIGHT = 64;

interface IPicAndTextProps {
  // Null or missing: no picture uploaded.
  imageSrc?: string | null;
  // Where the picture and the title link.
  url: string;
  // Only this part of the caption is the link: a name.
  title: string;
  // Plain text right after the link, e.g. " (12)".
  titleSuffix?: string;
}
export const PicAndText: React.FC<IPicAndTextProps> = ({ imageSrc, url, title, titleSuffix }) => {
  return (
    <Stack
      alignItems={"center"}
      flexDirection={"column"}
      sx={{ textAlign: "center" }}
    >
      <Box>
        <Link to={url}>
          <img height={IMAGE_HEIGHT} src={imageSrc ?? undefined} alt="image not found" />
        </Link>
      </Box>
      <Typography variant={"body2"}><Link to={url}>{title}</Link>{titleSuffix}</Typography>
    </Stack>
  );
};

// The loading state of a PicAndText. The square sits on the text baseline
// like the inline <img> does, so it leaves the same small gap under it.
export const PicAndTextSkeleton: React.FC = () => (
  <Stack alignItems={"center"} flexDirection={"column"}>
    <Box>
      <Skeleton
        variant={"rectangular"}
        width={IMAGE_HEIGHT}
        height={IMAGE_HEIGHT}
        sx={{ display: "inline-block" }}
      />
    </Box>
    <Typography variant={"body2"}><Skeleton variant={"text"} width={72} /></Typography>
  </Stack>
);
