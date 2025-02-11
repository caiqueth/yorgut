import { Box, Stack, Typography } from "@mui/material";

export const ProfileHeader = () => {
  return (
    <Box>
      <Stack direction={"row"}>
        <Typography>Nome Sobrenome</Typography>
      </Stack>
      <Stack direction={"row"}>
        <Typography>123 recados</Typography>
        <Typography>12 fotos</Typography>
        <Typography>1 vídeo</Typography>
        <Typography>8 fãs</Typography>
      </Stack>
      <Stack direction={"row"}>
        <Typography>Legal</Typography>
        <Typography>Confiável</Typography>
        <Typography>Lindo(a)</Typography>
      </Stack>
    </Box>
  );
};
