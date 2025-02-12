import { Box, Paper, Stack, Typography } from "@mui/material";

export const ProfileHeader = () => {
  const user = {
    name: "Frizz Chapper",
    statistics: {
      scraps: 123,
      photos: 12,
      videos: 1,
      fans: 8,
      trustworthy: 1,
      cool: 1,
      beautiful: 1,
    },
  };

  return (
    <Box>
      <Paper sx={{ backgroundColor: "white", p: 1 }}>
        <Stack direction={"row"} sx={{ width: "100%", padding: 0 }} alignItems={"end"}>
          <Typography variant={"h5"} sx={{ "&::first-letter": { fontSize: "200%" } }}>
            {user.name}
          </Typography>
        </Stack>
        <Stack direction={"row"} gap={1}>
          <Typography variant={"body2"}>📝 {user.statistics.scraps} recados</Typography>
          <Typography variant={"body2"}>📷 {user.statistics.photos} fotos</Typography>
          <Typography variant={"body2"}>📹 {user.statistics.videos} vídeo</Typography>
          <Typography variant={"body2"}>⭐ {user.statistics.fans} fãs</Typography>
        </Stack>
        <Stack direction={"row"} gap={1}>
          <Typography variant={"body2"}>{user.statistics.cool} Legal</Typography>
          <Typography variant={"body2"}>{user.statistics.trustworthy} Confiável</Typography>
          <Typography variant={"body2"}>{user.statistics.beautiful} Lindo(a)</Typography>
        </Stack>
      </Paper>
    </Box>
  );
};
