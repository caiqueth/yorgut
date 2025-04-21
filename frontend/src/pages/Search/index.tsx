import { Box, Stack, TextField, Typography } from "@mui/material";
import { Button } from "components/atoms/Button";
import React from "react";

interface IProps {}
export const Search: React.FC<IProps> = () => {
  return (
    <Box sx={{ width: "100%" }}>
      <Typography>Pesquisar por usuários ou comunidades</Typography>
      <Stack
        direction={"row"}
        alignItems={"center"}
        gap={1}
        sx={{ width: "100%" }}
      >
        <TextField />
        <Button>Pesquisar</Button>
      </Stack>
    </Box>
  );
};
