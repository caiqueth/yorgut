import React from "react";
import { useUser } from "api/user/queries";
import profileTextDecoration from "assets/images/profile_text_decoration.svg"

import { Box, Paper, Skeleton, Stack, Typography } from "@mui/material";
import { Repeater } from "components/atoms/Repeater";
import { useUserReviews } from "api/userReview/queries";
import { useUserFans } from "api/userFan/queries";
import { usePhotos } from "api/photo/queries";
import { useScraps } from "api/scrap/queries";

interface IProfileHeaderProps {
  userId: number;
}
export const ProfileHeader: React.FC<IProfileHeaderProps> = ({ userId }) => {
  const { data: user, isLoading } = useUser(userId);
  const { data: reviews } = useUserReviews(userId);
  const { data: fans } = useUserFans(userId);
  const { data: photos } = usePhotos(userId);
  // Same query as the scrapbook page, so both share one request.
  const { data: scraps } = useScraps(userId);

  const userStatistics = {
    totalScraps: scraps?.scraps.length ?? 0,
    totalPhotos: photos?.count ?? 0,
    totalVideos: 0,
    totalFans: fans?.fanCount ?? 0,
    trustworthy: Math.round(reviews?.trustworthy ?? 0),
    cool: Math.round(reviews?.cool ?? 0),
    sexy: Math.round(reviews?.sexy ?? 0),
  }

  return (
    <Box>
      <Paper sx={{ backgroundColor: "white", p: 1 }}>
        {isLoading
          ? <Skeleton height={116} />
          : (
            <Box>
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
                    {`${user?.firstName} ${user?.lastName}`}
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
            </Box>
          )
        }
      </Paper>
    </Box>
  );
};
