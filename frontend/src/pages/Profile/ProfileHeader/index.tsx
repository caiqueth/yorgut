import React from "react";
import { useUser } from "api/user/queries";
import profileTextDecoration from "assets/images/profile_text_decoration.svg"

import { Box, Paper, Stack, Typography } from "@mui/material";
import { Repeater } from "components/atoms/Repeater";

interface IProfileHeaderProps {
  userId: number;
}
export const ProfileHeader: React.FC<IProfileHeaderProps> = ({ userId }) => {
  const { data: user } = useUser(userId);

  const userStatistics = {
    totalScraps: 0,
    totalPhotos: 0,
    totalVideos: 0,
    totalFans: 0,
    trustworthy: Math.floor(10 * Math.random()),
    cool: Math.floor(10 * Math.random()),
    sexy: Math.floor(10 * Math.random()),
  }

  return (
    <Box>
      <Paper sx={{ backgroundColor: "white", p: 1 }}>
        <Stack direction={"row"} sx={{ width: "100%", padding: 0 }} alignItems={"end"}>
          <Box sx={{
            display: "flex",
          }}>
            <Box sx={{
              width: 45,
              height: 45,
              backgroundImage: `url(${profileTextDecoration})`,
              backgroundSize: "cover",
              backgroundRepeat: "no-repeat",
            }}>
            </Box>
            <Typography
              variant={"h5"}
              sx={{
                "&::first-letter": {
                  fontSize: 72,
                  ml: -3
                },
                lineHeight: 1
              }}>
              {user?.firstName} {user?.lastName}
            </Typography>
          </Box>

        </Stack>
        <Stack direction={"row"} spacing={1}>
          <Typography variant={"body2"}>📝 {userStatistics.totalScraps} recados</Typography>
          <Typography variant={"body2"}>📷 {userStatistics.totalPhotos} fotos</Typography>
          <Typography variant={"body2"}>📹 {userStatistics.totalVideos} vídeo</Typography>
          <Typography variant={"body2"}>⭐ {userStatistics.totalFans} fãs</Typography>
        </Stack>
        <Stack direction={"row"} spacing={0}>
          <Repeater element={"🧊"} times={userStatistics.cool} />
          <Repeater element={"😀"} times={userStatistics.trustworthy} />
          <Repeater element={"❤️"} times={userStatistics.sexy} />
        </Stack>
      </Paper>
    </Box>
  );
};
