import React from "react";
import { Stack, Typography } from "@mui/material";
import { Link } from "react-router-dom";

interface IProps {
  imageSrc: string;
  url: string;
  text: string;
}
export const PicAndText: React.FC<IProps> = ({ imageSrc, url, text }) => {
  return (
    <Stack
      alignItems={"center"}
      flexDirection={"column"}
      sx={{ textAlign: "center" }}
    >
      <Link to={url}>
        <img height={64} src={imageSrc} alt="image not found" />
        <Typography variant={"body2"}>{text}</Typography>
      </Link>
    </Stack>
  );
};
