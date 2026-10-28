import { Typography } from "@mui/material";
import React from "react";
import { Link } from "react-router-dom";

interface IProps {
  // In pixels, or by breakpoint ({ xs: 22, sm: 36 }).
  size?: number | Partial<Record<"xs" | "sm" | "md" | "lg" | "xl", number>>;
}
// The logo, linking to the home page.
export const YorgutTitle: React.FC<IProps> = ({ size }) => {
  return <Typography
    variant={"h5"}
    component={Link}
    to={"/"}
    sx={{
      // Doubled so it beats the navbar's and the theme's link colors and
      // underline ("& a", "a:hover").
      "&&, &&:visited, &&:hover": {
        color: "primary.main",
        textDecoration: "none",
      },
      fontFamily: "'Quicksand Variable', sans-serif",
      fontWeight: 800,
      fontSize: size ?? "inherit",
      lineHeight: 1
    }}
  >yorgut
  </Typography>
}
