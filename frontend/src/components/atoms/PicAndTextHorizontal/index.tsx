import React from "react";
import { Box, Stack, Typography } from "@mui/material";
import { Link } from "react-router-dom";

interface IProps {
  imageSrc: string;
  url: string;
  title: string;
  text?: string | null;
  imageSize?: number;
}
export const PicAndTextHorizontal: React.FC<IProps> = ({
  url,
  title,
  text,
  imageSrc,
  imageSize = 40
}) => {
  return (
    <Stack direction={"row"} gap={1} alignItems={"center"}>
      <Box sx={{
        display: "block",
        minWidth: imageSize,
        width: imageSize,
        height: imageSize,
        minHeight: imageSize,
        backgroundImage: `url("${imageSrc}")`,
        backgroundSize: "cover",
        backgroundRepeat: "no-repeat",
        alignSelf: "start",
      }} />
      <Stack direction={"column"}>
        <Typography><Link to={url}>{title}</Link></Typography>
        {text && <Typography sx={{
          whiteSpace: "pre-line"
        }}>{text}</Typography>}
      </Stack>
    </Stack>
  );
};
