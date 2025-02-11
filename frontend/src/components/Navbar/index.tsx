import React from "react";

import { Stack, Typography } from "@mui/material";

export const Navbar: React.FC = () => {
  return (
    <Stack direction={"row"}>
      <Typography>Home</Typography>
      <Typography>Amigos</Typography>
      <Typography>Mensagens</Typography>
      <Typography>Comunidades</Typography>
      <Typography>Pesquisar</Typography>
      <Typography>Mídia</Typography>
      <Typography>O que há de novo</Typography>
    </Stack>
  );
};
