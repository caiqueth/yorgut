import { Typography } from "@mui/material";
import React from "react";

interface IProps {
  size?: number;
}
export const YorgutTitle: React.FC<IProps> = ({ size }) => {
  return <Typography
    variant={"h5"}
    sx={{
      color: "#d744b6",
      fontFamily: "Quicksand",
      fontWeight: 800,
      fontSize: size ?? "inherit",
      lineHeight: 1
    }}
  >yorgut
  </Typography>
}