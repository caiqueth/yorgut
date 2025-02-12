import { Box, Paper, Stack, Typography } from "@mui/material";

export const ProfileHeader = () => {
  return (
    <Box>
      <Paper sx={{ backgroundColor: "white", p: 1 }}>
        <Stack direction={"row"} sx={{ width: "100%", padding: "22px 0 0 6px" }}>
          <Typography variant={"h5"}>Frizz Chapper</Typography>
        </Stack>
        <Stack direction={"row"} gap={1}>
          <Typography>📝 123 recados</Typography>
          <Typography>📷 12 fotos</Typography>
          <Typography>📹 1 vídeo</Typography>
          <Typography>⭐ 8 fãs</Typography>
        </Stack>
        <Stack direction={"row"} gap={1}>
          <Typography>Legal</Typography>
          <Typography>Confiável</Typography>
          <Typography>Lindo(a)</Typography>
        </Stack>
      </Paper>
    </Box>
  );
};
