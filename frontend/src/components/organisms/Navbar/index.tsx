import React from "react";
import { Link } from "react-router-dom";
import { palette } from "theme";

import { Box, Stack, Typography } from "@mui/material";

export const Navbar: React.FC = () => {
  return (
    <Box sx={{ width: "100%", backgroundColor: palette.darkBlue, mb: 2 }}>
      <Stack direction={"row"} gap={1} sx={{ width: "100%" }}>
        <Typography variant={"body2"}>
          <Link to="/">Home</Link>
        </Typography>
        <Typography variant={"body2"}>|</Typography>
        <Typography variant={"body2"}>
          <Link to="/friends">Amigos</Link>
        </Typography>
        <Typography variant={"body2"}>|</Typography>
        <Typography variant={"body2"}>
          <Link to="/scraps">Mensagens</Link>
        </Typography>
        <Typography variant={"body2"}>|</Typography>
        <Typography variant={"body2"}>
          <Link to="/communities">Comunidades</Link>
        </Typography>
        <Typography variant={"body2"}>|</Typography>
        <Typography variant={"body2"}>
          <Link to="/search">Pesquisar</Link>
        </Typography>
        <Typography variant={"body2"}>|</Typography>
        <Typography variant={"body2"}>
          <Link to="/media">Mídia</Link>
        </Typography>
        <Typography variant={"body2"}>|</Typography>
        <Typography variant={"body2"}>
          <Link to="/whats-new">O que há de novo</Link>
        </Typography>
      </Stack>
    </Box>
  );
};
