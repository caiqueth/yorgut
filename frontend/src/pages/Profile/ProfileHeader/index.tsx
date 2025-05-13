import React from "react";
import { useUser, useUserStatistics } from "api/user/queries";

import { Box, Paper, Stack, Typography } from "@mui/material";

interface IProfileHeaderProps {
  userId: number;
}
export const ProfileHeader: React.FC<IProfileHeaderProps> = ({ userId }) => {
  const { data: user } = useUser(userId);
  const userStatistics = useUserStatistics();

  return (
    <Box>
      <Paper sx={{ backgroundColor: "white", p: 1 }}>
        <Stack direction={"row"} sx={{ width: "100%", padding: 0 }} alignItems={"end"}>
          <Typography variant={"h5"} sx={{ "&::first-letter": { fontSize: "200%" } }}>
            {user?.firstName} {user?.lastName}
          </Typography>
        </Stack>
        <Stack direction={"row"} gap={1}>
          <Typography variant={"body2"}>📝 {userStatistics.totalScraps} recados</Typography>
          <Typography variant={"body2"}>📷 {userStatistics.totalPhotos} fotos</Typography>
          <Typography variant={"body2"}>📹 {userStatistics.totalVideos} vídeo</Typography>
          <Typography variant={"body2"}>⭐ {userStatistics.totalFans} fãs</Typography>
        </Stack>
        <Stack direction={"row"} spacing={2}>
          <Typography variant={"body2"} sx={{ letterSpacing: -12 }}>🧊🧊🧊🧊🧊🧊🧊🧊🧊🧊</Typography>
          <Typography variant={"body2"} sx={{ letterSpacing: -12 }}>😀😀😀😀😀😀😀😀😀😀</Typography>
          <Typography variant={"body2"} sx={{ letterSpacing: -12 }}>❤️❤️❤️❤️❤️❤️❤️❤️❤️❤️</Typography>
        </Stack>
      </Paper>
    </Box>
  );
};
