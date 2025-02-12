import React from "react";
import { palette } from "theme";

import { Box, Stack, Typography } from "@mui/material";

export const Navbar: React.FC = () => {
  return (
    <Box sx={{ width: "100%", backgroundColor: palette.darkBlue, mb: 2 }}>
      <Stack direction={"row"} gap={1} sx={{ width: "100%" }}>
        <Typography variant={"body2"}>Home</Typography>
        <Typography variant={"body2"}>|</Typography>
        <Typography variant={"body2"}>Amigos</Typography>
        <Typography variant={"body2"}>|</Typography>
        <Typography variant={"body2"}>Mensagens</Typography>
        <Typography variant={"body2"}>|</Typography>
        <Typography variant={"body2"}>Comunidades</Typography>
        <Typography variant={"body2"}>|</Typography>
        <Typography variant={"body2"}>Pesquisar</Typography>
        <Typography variant={"body2"}>|</Typography>
        <Typography variant={"body2"}>Mídia</Typography>
        <Typography variant={"body2"}>|</Typography>
        <Typography variant={"body2"}>O que há de novo</Typography>
      </Stack>
    </Box>
  );
};
