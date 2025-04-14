import { Box, Paper, Typography } from "@mui/material";
import React from "react";
import { palette } from "theme";

interface IBoxWithTitleProps {
  title: string;
  children?: React.ReactNode;
}
export const BoxWithTitle: React.FC<IBoxWithTitleProps> = ({
  title,
  children,
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
      {children}
    </Paper>
  );
};
