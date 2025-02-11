import React from "react";

import { Box } from "@mui/material";

interface IWrapperProps {
  children: React.ReactNode;
}
export const Wrapper: React.FC<IWrapperProps> = ({ children }) => {
  return <Box sx={{ width: "1024px", display: "flex", alignContent: "center" }}>{children}</Box>;
};
