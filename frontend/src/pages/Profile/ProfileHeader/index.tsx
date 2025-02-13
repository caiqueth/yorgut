import React from "react";
import { useUserDetails, useUserStatistics } from "api/user/queries";

import { Box, Paper, Stack, Typography } from "@mui/material";

interface IProfileHeaderProps {
  userId: number;
}
export const ProfileHeader: React.FC<IProfileHeaderProps> = ({ userId }) => {
  const user = useUserDetails(userId);
  const userStatistics = useUserStatistics(userId);

  return (
    <Box>
      <Paper sx={{ backgroundColor: "white", p: 1 }}>
        <Stack direction={"row"} sx={{ width: "100%", padding: 0 }} alignItems={"end"}>
          <Typography variant={"h5"} sx={{ "&::first-letter": { fontSize: "200%" } }}>
            {user.firstName} {user.lastName}
          </Typography>
        </Stack>
        <Stack direction={"row"} gap={1}>
          <Typography variant={"body2"}>📝 {userStatistics.totalScraps} recados</Typography>
          <Typography variant={"body2"}>📷 {userStatistics.totalPhotos} fotos</Typography>
          <Typography variant={"body2"}>📹 {userStatistics.totalVideos} vídeo</Typography>
          <Typography variant={"body2"}>⭐ {userStatistics.totalFans} fãs</Typography>
        </Stack>
        <Stack direction={"row"} gap={1}>
          <Typography variant={"body2"}>{userStatistics.cool} Legal</Typography>
          <Typography variant={"body2"}>{userStatistics.trustworthy} Confiável</Typography>
          <Typography variant={"body2"}>{userStatistics.beautiful} Lindo(a)</Typography>
        </Stack>
      </Paper>
    </Box>
  );
};
