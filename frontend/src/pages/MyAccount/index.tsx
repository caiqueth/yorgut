import { Box, Typography } from "@mui/material";
import React from "react";

interface IProps {}
export const MyAccount: React.FC<IProps> = () => {
  return (
    <Box>
      <Typography variant={"body1"}>
        Aqui vão opções gerais para gerenciamento da conta
      </Typography>
    </Box>
  );
};
