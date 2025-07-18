import { Box, LinearProgress, Paper, Typography } from "@mui/material";
import React from "react";
import { palette } from "theme";

interface IBoxWithTitleProps {
  title: string;
  children?: React.ReactNode;
  isLoading?: boolean;
}
export const BoxWithTitle: React.FC<IBoxWithTitleProps> = ({
  title,
  children,
  isLoading,
}) => {
  return (
    <Paper>
      <Box
        sx={{
          textAlign: "center",
          width: "100%",
          backgroundColor: palette.lightBlue,
        }}
      >
        <Typography variant={"subtitle2"} fontWeight={"bold"}>
          {title}
        </Typography>
      </Box>
      <LinearProgress sx={{ visibility: isLoading ? "visible" : "hidden" }} />
      {children}
    </Paper>
  );
};
